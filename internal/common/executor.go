package common

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/nuccitheboss/gambol/internal/provider/lxd"
	"github.com/nuccitheboss/gambol/internal/storage"
)

var (
	provider lxd.Driver
	cache    storage.Cache
)

// Run gambol playthrough.
func Run(file string) error {
	play, err := loadPlay(file)
	if err != nil {
		return err
	}

	queue := assembleWorkQueue(play)

	provider, err = lxd.New()
	if err != nil {
		return err
	}

	cache, err = storage.NewCache(uuid.NewString())
	if err != nil {
		return err
	}

	if err := runActs(queue); err != nil {
		return err
	}

	ids, err := cache.GetInstanceIds()
	if err != nil {
		return err
	}
	if err := provider.DestroyInstances(ids); err != nil {
		return err
	}
	cache.Flush()

	return nil
}

func runActs(queue WorkQueue) error {
	for !queue.IsEmpty() {
		act, err := queue.Pop()
		if err != nil {
			break
		}
		if err := runAct(act); err != nil {
			return err
		}
	}

	return nil
}

func runAct(act Act) error {
	fmt.Printf("Executing Act: %s\n", act.Option.Name)
	exists, err := provider.CheckIfInstanceExists(act.Option.RunOn)
	if err != nil {
		return err
	}

	var instanceId string
	if !exists {
		if err := provider.CreateInstance(act.Id, act.Option.RunOn); err != nil {
			return err
		}
		if err := cache.PutInstanceId(act.Id); err != nil {
			return err
		}
		instanceId = act.Id
	} else {
		instanceId = act.Option.RunOn
	}

	if len(act.Option.Input) > 0 {
		if err := push(instanceId, act.Option.Input); err != nil {
			return err
		}
	}

	if err := runScenes(instanceId, act.Option.Scenes); err != nil {
		return err
	}

	if len(act.Option.Output) > 0 {
		if err := pull(instanceId, act.Option.Output); err != nil {
			return err
		}
	}

	if !act.Option.KeepAlive {
		if err := provider.StopInstance(instanceId); err != nil {
			return err
		}
	}

	return nil
}

func runScenes(id string, scenes []Scene) error {
	for _, scene := range scenes {
		fmt.Printf("Executing Scene: %s\n", scene.Name)
		if err := provider.ExecInstance(id, scene.Run); err != nil {
			return err
		}
	}

	return nil
}
