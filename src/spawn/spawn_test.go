package spawn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpawn(t *testing.T) {
	assert := assert.New(t)

	cmds := []string{
		"sleep 1",
		"sleep 1",
		"sleep 1",
		"sleep 1",
	}

	spawner := NewSpawner(2)
	spawner.Add(cmds...)
	spawner.Run()
	assert.Nil(spawner.Err())

	spawner = NewSpawner(0)
	spawner.Add(cmds...)
	spawner.Run()
	assert.Nil(spawner.Err())

	err := Spawn("sleep 1")
	assert.Nil(err)
}
