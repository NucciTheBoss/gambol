package storage

import (
	"os"

	"github.com/spf13/viper"
)

func Init() error {
	storagePath := viper.GetString("storage")
	err := os.MkdirAll(storagePath, 0700)
	if err != nil {
		return err
	}

	return nil
}
