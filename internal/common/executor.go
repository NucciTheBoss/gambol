package common

import (
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
	exists, err := provider.CheckIfInstanceExists(act.Id)
	if err != nil {
		return err
	} else if !exists {
		if err := provider.CreateInstance(act.Id, act.Option.RunOn); err != nil {
			return err
		}
		if err := cache.PutInstanceId(act.Id); err != nil {
			return err
		}
	}

	if len(act.Option.Input) > 0 {
		if err := push(act.Id, act.Option.Input); err != nil {
			return err
		}
	}

	if err := runScenes(act.Id, act.Option.Scenes); err != nil {
		return err
	}

	if len(act.Option.Output) > 0 {
		if err := pull(act.Id, act.Option.Output); err != nil {
			return err
		}
	}

	if !act.Option.KeepAlive {
		if err := provider.StopInstance(act.Id); err != nil {
			return err
		}
	}

	return nil
}

func runScenes(id string, scenes []Scene) error {
	for _, scene := range scenes {
		if err := provider.ExecInstance(id, scene.Run); err != nil {
			return err
		}
	}

	return nil
}
