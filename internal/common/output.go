package common

import (
	"github.com/nuccitheboss/gambol/internal/storage"
)

// Pull artifacts from an Act instance.
func pull(id string, artifacts []storage.Artifact) error {
	for _, artifact := range artifacts {
		if artifact.HostPath != "" {
			if err := pullHostPath(id, artifact); err != nil {
				return err
			}
		} else {
			if err := pullCache(id, artifact); err != nil {
				return err
			}
		}
	}

	return nil
}

// Pull artifact from an Act instance and store in playthrough cache.
func pullCache(id string, artifact storage.Artifact) error {
	output, err := provider.GetArtifact(id, artifact)
	if err != nil {
		return err
	}
	if err := cache.PutArtifact(artifact.Key, output); err != nil {
		return err
	}

	return nil
}

// Pull artifact from an Act instance and unwrap on host.
func pullHostPath(id string, artifact storage.Artifact) error {
	output, err := provider.GetArtifact(id, artifact)
	if err != nil {
		return err
	}
	if err := artifact.Unwrap(output); err != nil {
		return err
	}

	return nil
}
