package storage

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// `Artifact` represents a possible input or output from an Act run.
// Inputs will be pushed into the Act instance before any scenes are executed,
// and Outputs will be cached locally after all steps have completed successfully.
type Artifact struct {
	// Key for storing or retrieving artifact from cache.
	Key string `yaml:"key"`

	// Path to retrieve or dump artifact on host.
	HostPath string `yaml:"host-path"`

	// Path to push or pull artifact to or from
	// in Act instance.
	Path string `yaml:"path"`
}

// Get unique name of the artifact. If both a `key and `host-path` are
// provided in the playthrough file, then an error is returned.
func (a *Artifact) Name() (string, error) {
	if a.HostPath != "" && a.Key != "" {
		return "", fmt.Errorf("artifact can only have a unique host-path or key")
	} else if a.Key != "" {
		return a.Key, nil
	} else {
		return a.HostPath, nil
	}
}

func (a *Artifact) IsDir(artifact []byte) (bool, error) {
	r := bytes.NewBuffer(artifact)
	tr := tar.NewReader(r)

	name, err := a.Name()
	if err != nil {
		return false, err
	}

	// Peek into tar file to determine if artifact
	// is a directory or a file.
	hdr, err := tr.Next()
	if err == io.EOF {
		return false, fmt.Errorf("artifact '%s' is empty", name)
	}
	if err != nil {
		return false, err
	}

	return hdr.FileInfo().IsDir(), nil
}

// Wrap an artifact inside a tar archive pulled from the configured `host-path`.
// An error is returned if an error occurs when wrapping the artifact pulled
// from the given location on the host. Returns an artifact wrapped in a tarball.
func (a *Artifact) Wrap() ([]byte, error) {
	var artifact bytes.Buffer
	tw := tar.NewWriter(&artifact)

	info, err := os.Stat(a.HostPath)
	if err != nil {
		return nil, err
	}

	var base string
	if info.IsDir() {
		base = filepath.Base(a.HostPath) + "/"
	}
	if err := filepath.Walk(a.HostPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		if base != "" {
			hdr.Name = filepath.Join(base, strings.TrimPrefix(path, a.HostPath))
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		fin, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fin.Close()
		_, err = io.Copy(tw, fin)
		return err
	}); err != nil {
		return nil, err
	}

	tw.Close()
	return io.ReadAll(&artifact)
}

// Unwrap a tarred artifact and dump to the configured `host-path`.
// An error will be returned if an error occurs when unwrapping the
// artifact and dumping to the given location on the host.
func (a *Artifact) Unwrap(artifact []byte) error {
	r := bytes.NewReader(artifact)
	tr := tar.NewReader(r)

	dir, err := a.IsDir(artifact)
	if err != nil {
		return err
	}

	var root string
	if dir {
		if err := os.MkdirAll(a.HostPath, 0755); err != nil {
			return err
		}
	}
	for {
		var path string
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		info := hdr.FileInfo()

		if dir {
			if root == "" {
				root, _ = filepath.Split(hdr.Name)
			}
			strippedName, _ := filepath.Rel(root, hdr.Name)
			if strippedName == "." {
				continue
			}
			path = filepath.Join(a.HostPath, strippedName)
		} else {
			root = strings.TrimSuffix(a.HostPath, filepath.Base(a.HostPath))
			if _, err := os.Stat(root); err != nil {
				if err := os.MkdirAll(root, 0755); err != nil {
					return err
				}
			}
			path = a.HostPath
		}

		if info.IsDir() {
			if err := os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}

			continue
		}

		fout, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer fout.Close()
		if _, err := io.Copy(fout, tr); err != nil {
			return err
		}
	}

	return nil
}
