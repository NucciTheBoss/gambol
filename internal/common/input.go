package common

import (
	"github.com/nuccitheboss/gambol/internal/storage"
)

// Push artifacts into an Act instance.
func push(id string, artifacts []storage.Artifact) error {
	for _, artifact := range artifacts {
		if artifact.HostPath != "" {
			if err := pushHostPath(id, artifact); err != nil {
				return err
			}
		} else {
			if err := pushCache(id, artifact); err != nil {
				return err
			}
		}
	}

	return nil
}

// Push artifact located in the playthrough cache into Act instance.
func pushCache(id string, artifact storage.Artifact) error {
	input, err := cache.GetArtifact(artifact.Key)
	if err != nil {
		return err
	}
	if err := provider.PutArtifact(id, artifact, input); err != nil {
		return err
	}

	return nil
}

// Push artifact located on host into Act instance.
func pushHostPath(id string, artifact storage.Artifact) error {
	input, err := artifact.Wrap()
	if err != nil {
		return err
	}
	if err := provider.PutArtifact(id, artifact, input); err != nil {
		return err
	}

	return nil
}
