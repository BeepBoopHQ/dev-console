package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBotRegistryRemove(t *testing.T) {
	assert := assert.New(t)
	msg := Message{
		Type:        "add_resource",
		Date:        time.Now(),
		MsgID:       "MsgID",
		BotResource: BotResource{ResourceID: "resID", Resource: map[string]string{"howdy": "doody"}},
	}

	registry.Add(msg)
	assert.Equal("resID", registry.botResources[0].ResourceID)
}
