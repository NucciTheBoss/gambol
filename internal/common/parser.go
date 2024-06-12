package common

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Load gambol playthrough file.
func loadPlay(file string) (play Play, err error) {
	runData, err := os.ReadFile(file)
	if err != nil {
		return play, err
	}

	err = yaml.Unmarshal(runData, &play)
	if err != nil {
		return play, err
	}

	return play, nil
}

// Assemble the work queue for the Act executor.
//
// FIXME: Work queue assembly is quite rudementary at the moment,
// but it will get fancier as I add parallelism and more
// user control over how Act execution is scheduled.
// Right now it's just synchronous because that's the easiest ;).
func assembleWorkQueue(play Play) (w WorkQueue) {
	w = new(WorkQueue).Init(play.Acts)
	return w
}
