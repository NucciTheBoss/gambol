package common

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/nuccitheboss/gambol/internal/storage"
)

// `Play` is the top-level data structure that represents
// a single gambol playthrough.
type Play struct {
	// Name of the playthrough.
	Name string `yaml:"name"`

	// Provider to use for providing the instances
	// that act scenes will be run within.
	Provider map[string]Provider `yaml:"provider"`

	// Acts of the playthrough. Each act is mapped
	// to a single instance that scenes will be run within.
	Acts Acts `yaml:"acts"`
	// Acts map[string]Act `yaml:"acts"`
}

// `Provider` represents an Act instance provider.
// Providers are used for providing the isolated instances
// that Acts and their complementing scenes are executed within.
type Provider struct {
}

type Acts []Act

func (a *Acts) UnmarshalYAML(v *yaml.Node) error {
	if v.Kind != yaml.MappingNode {
		return fmt.Errorf("`acts` must be a valid YAML mapping, not %v", v.Kind)
	}

	*a = make([]Act, len(v.Content)/2)
	for i := 0; i < len(v.Content); i += 2 {
		act := &(*a)[i/2]
		err := v.Content[i].Decode(&act.Id)
		if err != nil {
			return err
		}
		err = v.Content[i+1].Decode(&act.Option)
		if err != nil {
			return err
		}
	}

	return nil
}

// `Act` is a job to run within a single instance.
// An instance will be provisioned by the provider for each Act.
// Acts can also be executed within the same instance.
type Act struct {
	// Unique id of Act.
	Id string

	// Options that will processed by executor.
	Option ActOptions
}

type ActOptions struct {
	// Name of Act.
	Name string `yaml:"name"`

	// Base image to run the act scenes within.
	// The base image name must map to a valid
	// image name within the instance instance provider.
	// `On` can also be mapped to an already provisioned
	// instance to execute act scenes within the same instance.
	RunOn string `yaml:"run-on"`

	// If true, keep instance running after Act execution has
	// completed. Otherwise, shut down instance to free up resources
	// for other acts. Useful for distributed systems testing.
	KeepAlive bool `yaml:"keep-alive"`

	// Input data to push into act instance before executing scenes.
	Input []storage.Artifact `yaml:"input"`

	// Output data to pull from act instance after scenes have been executed.
	Output []storage.Artifact `yaml:"output"`

	// Act scenes. Each scene will be executed
	// in sequential order within the act instance.
	Scenes []Scene `yaml:"scenes"`
}

// `Scene` is a step within an Act. Each Scene is
// executed within the same Act instance in sequential order.
type Scene struct {
	// Name of the scene.
	Name string `yaml:"name"`

	// Run script to execute within the act instance.
	Run string `yaml:"run"`
}

// Queue of Acts to be handled by executor.
type WorkQueue struct {
	acts []Act
}

func (w WorkQueue) Init(acts []Act) WorkQueue {
	w.acts = make([]Act, len(acts))
	copy(w.acts, acts)
	return w
}

// Pop Act from the front of the work queue.
func (w *WorkQueue) Pop() (act Act, err error) {
	if len(w.acts) == 0 {
		return act, errors.New("no acts in work queue")
	}

	act = w.acts[0]
	if len(w.acts) == 1 {
		w.acts = nil
	} else {
		w.acts = w.acts[1:]
	}

	return act, nil
}

func (w *WorkQueue) IsEmpty() bool {
	return len(w.acts) == 0
}
