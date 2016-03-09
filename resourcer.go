package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pborman/uuid"
)

type BotResourceRegistry struct {
	BotResources []Resource `json:"resources"`
}

// a global for now
var registry = &BotResourceRegistry{}

func (brr *BotResourceRegistry) String() string {
	var str string
	for _, item := range brr.BotResources {
		str += fmt.Sprintf("%+v\n", item)
	}
	return str
}

func (brr *BotResourceRegistry) Json() string {
	data, _ := json.Marshal(brr.BotResources)
	return string(data)
}

func (brr *BotResourceRegistry) Add(resourceID string, resource *Resource) {
	addResMsg := NewAddResourceMessage(resourceID, resource)

	if wsConn != nil {
		wsConn.WriteJSON(addResMsg)
	}

	resourceContainer := Resource{
		ID:          resourceID,
		ResourceCfg: resource.ResourceCfg,
		UserID:      "unknown",
		PrototypeID: idGen(),
		Time:        time.Now(),
	}
	brr.BotResources = append(brr.BotResources, resourceContainer)
}

func (brr *BotResourceRegistry) Update(msg *UpdateResourceMessage) {
	for i, registryItem := range brr.BotResources {
		if registryItem.ID == msg.ResourceID {
			brr.BotResources[i] = Resource{ID: msg.ResourceID, ResourceCfg: msg.Resource}
		}
	}
}

func (brr *BotResourceRegistry) Remove(msg *RemoveResourceMessage) {
	for i, registryItem := range brr.BotResources {
		if registryItem.ID == msg.ResourceID {
			brr.BotResources = append(brr.BotResources[:i], brr.BotResources[i+1:]...)
		}
	}
}

func idGen() string {
	return uuid.NewRandom().String()
}
