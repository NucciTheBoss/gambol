package lxd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
	"github.com/google/uuid"
	"github.com/nuccitheboss/gambol/internal/storage"
)

type Driver struct {
	server lxd.InstanceServer
}

// Establish connection to LXD server.
//
// FIXME: Currently just using the Unix socket
// because that's the easiest. Will get more advanced
// later once further requirements and use cases are clear.
func New() (provider Driver, err error) {
	server, err := lxd.ConnectLXDUnix("/var/snap/lxd/common/lxd/unix.socket", nil)
	if err != nil {
		return provider, err
	}

	provider.server = server
	return provider, nil
}

func (p *Driver) CheckIfInstanceExists(id string) (bool, error) {
	names, err := p.server.GetInstanceNames("container")
	if err != nil {
		return false, err
	}

	return slices.Contains(names, id), nil
}

func (p *Driver) CheckIfInstanceActive(id string) (bool, error) {
	instance, _, err := p.server.GetInstanceState(id)
	if err != nil {
		return false, err
	}
	switch instance.StatusCode {
	case 102: // Stopped
		return false, nil
	case 112: // Error
		return false, nil
	default:
		return true, nil
	}
}

func (p *Driver) CreateInstance(id string, platform string) error {
	request := api.InstancesPost{
		Name: id,
		Source: api.InstanceSource{
			Type:     "image",
			Protocol: "simplestreams",
			Server:   "https://cloud-images.ubuntu.com/releases",
			Alias:    platform,
		},
		Type: "container",
	}
	op, err := p.server.CreateInstance(request)
	if err != nil {
		return err
	}
	err = op.Wait()
	if err != nil {
		return err
	}

	startRequest := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1,
	}
	op, err = p.server.UpdateInstanceState(id, startRequest, "")
	if err != nil {
		return err
	}
	err = op.Wait()
	if err != nil {
		return err
	}

	execRequest := api.InstanceExecPost{
		Command: []string{"mkdir", "-p", "/root/.gambol/input", "/root/.gambol/output"},
	}
	execArgs := lxd.InstanceExecArgs{}
	op, err = p.server.ExecInstance(id, execRequest, &execArgs)
	if err != nil {
		return err
	}
	err = op.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (p *Driver) StopInstance(id string) error {
	stopRequest := api.InstanceStatePut{
		Action:  "stop",
		Timeout: -1,
	}
	op, err := p.server.UpdateInstanceState(id, stopRequest, "")
	if err != nil {
		return err
	}

	if err := op.Wait(); err != nil {
		return err
	}

	return nil
}

func (p *Driver) DestroyInstances(ids []string) error {
	for _, id := range ids {
		active, err := p.CheckIfInstanceActive(id)
		if err != nil {
			return nil
		}
		if active {
			if err := p.StopInstance(id); err != nil {
				return err
			}
		}
		if _, err := p.server.DeleteInstance(id); err != nil {
			return err
		}
	}

	return nil
}

func (p *Driver) ExecInstance(id string, script string) error {
	uploadArgs := lxd.InstanceFileArgs{
		Content:   bytes.NewReader([]byte(script)),
		WriteMode: "overwrite",
		Type:      "file",
	}
	err := p.server.CreateInstanceFile(id, "/root/.gambol/run", uploadArgs)
	if err != nil {
		return err
	}

	execRequest := api.InstanceExecPost{
		Command:      []string{"bash", "/root/.gambol/run"},
		RecordOutput: true,
	}
	execArgs := lxd.InstanceExecArgs{}
	op, err := p.server.ExecInstance(id, execRequest, &execArgs)
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}

	if returnCode := op.Get().Metadata["return"]; returnCode != float64(0) {
		return errors.New("scene failed")
	}

	return nil
}

var getArtifactScript = `
UUID=%s
TARGET=%s
pushd $(dirname ${TARGET})
tar -cf /root/.gambol/output/${UUID}.tar ${TARGET}
`

// Get Artifact from Act instance.
func (p *Driver) GetArtifact(id string, artifact storage.Artifact) (out []byte, err error) {
	uniqueID := uuid.NewString()
	wrapper := fmt.Sprintf(getArtifactScript, uniqueID, artifact.Path)
	if err := p.ExecInstance(id, wrapper); err != nil {
		return nil, err
	}

	target := fmt.Sprintf("/root/.gambol/output/%s.tar", uniqueID)
	buf, _, err := p.server.GetInstanceFile(id, target)
	if err != nil {
		return nil, err
	}

	out, err = io.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	return out, nil
}

var putFileArtifactScript = `
UUID=%s
TARGET=%s
OUTPUT=$(mktemp -p /tmp -d)
tar -xf /root/.gambol/input/${UUID}.tar -C ${OUTPUT}
mkdir -p $(dirname ${TARGET})
mv ${OUTPUT}/* ${TARGET}
`
var putDirArtifactScript = `
UUID=%s
TARGET=%s
mkdir -p ${TARGET}
tar -xf /root/.gambol/input/${UUID}.tar -C ${TARGET} --strip-components 1
`

// Put (upload) Artifact into Act instance.
func (p *Driver) PutArtifact(id string, artifact storage.Artifact, input []byte) error {
	uniqueID := uuid.NewString()
	target := fmt.Sprintf("/root/.gambol/input/%s.tar", uniqueID)
	args := lxd.InstanceFileArgs{
		Content:   bytes.NewReader(input),
		Type:      "file",
		WriteMode: "overwrite",
	}
	if err := p.server.CreateInstanceFile(id, target, args); err != nil {
		return err
	}

	dir, err := artifact.IsDir(input)
	if err != nil {
		return err
	}

	var wrapper string
	if !dir {
		wrapper = fmt.Sprintf(putFileArtifactScript, uniqueID, artifact.Path)
	} else {
		wrapper = fmt.Sprintf(putDirArtifactScript, uniqueID, artifact.Path)
	}
	if err := p.ExecInstance(id, wrapper); err != nil {
		return err
	}

	return nil
}
