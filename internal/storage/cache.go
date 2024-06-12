package storage

import (
	"errors"
	"fmt"
	"path"

	"github.com/spf13/viper"
	"go.etcd.io/bbolt"
)

type Cache struct {
	// Unique name for gambol playthrough cache buckets.
	name []byte

	// Path to artifacts cache.
	artifactDB string

	// Path to instance cache.
	instanceDB string
}

func NewCache(name string) (c Cache, err error) {
	c.name = []byte(name)
	storagePath := viper.GetString("storage")
	c.artifactDB = path.Join(storagePath, "artifact.db")
	c.instanceDB = path.Join(storagePath, "instance.db")
	err = c.init()
	if err != nil {
		return c, err
	}
	return c, nil
}

func (c *Cache) PutArtifact(key string, artifact []byte) error {
	db, err := c.openArtifactDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(c.name)
		if b == nil {
			return errors.New("failed to open artifact cache")
		}

		return b.Put([]byte(key), artifact)
	}); err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetArtifact(key string) ([]byte, error) {
	db, err := c.openArtifactDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var artifact []byte
	if err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(c.name)
		if b == nil {
			return errors.New("failed to open artifact cache")
		}
		v := b.Get([]byte(key))
		artifact = make([]byte, len(v))
		copy(artifact, v)

		return nil
	}); err != nil {
		return nil, err
	}

	if artifact == nil {
		return nil, fmt.Errorf("artifact '%s' not found in cache", key)
	}

	return artifact, nil
}

func (c *Cache) PutInstanceId(id string) {

}

func (c *Cache) GetInstanceIds() {

}

// Flush out caches after successful completion of playthrough.
func (c *Cache) Flush() error {
	adb, err := c.openArtifactDB()
	if err != nil {
		return err
	}
	defer adb.Close()

	if err := adb.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket(c.name)
	}); err != nil {
		return err
	}

	idb, err := c.openInstanceDB()
	if err != nil {
		return err
	}
	defer idb.Close()

	if err := idb.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket(c.name)
	}); err != nil {
		return err
	}

	return nil
}

// Initialize artifact and instance database.
func (c *Cache) init() error {
	adb, err := c.openArtifactDB()
	if err != nil {
		return err
	}
	defer adb.Close()

	err = adb.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket(c.name)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	idb, err := c.openInstanceDB()
	if err != nil {
		return err
	}
	defer idb.Close()

	err = idb.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket(c.name)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) openArtifactDB() (*bbolt.DB, error) {
	db, err := bbolt.Open(c.artifactDB, 0600, &bbolt.Options{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *Cache) openInstanceDB() (*bbolt.DB, error) {
	db, err := bbolt.Open(c.instanceDB, 0600, &bbolt.Options{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
