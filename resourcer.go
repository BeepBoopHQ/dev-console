package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/pborman/uuid"
)

type ConsolePostMessage struct {
	ResourceID string `json:"resourceID"`
	Resource
}

type BotResourceRegistry struct {
	botResources []Resource
}

var registry = &BotResourceRegistry{}

func (brr *BotResourceRegistry) String() string {
	var str string
	for _, item := range brr.botResources {
		str += fmt.Sprintf("%+v\n", item)
	}
	return str
}

func (brr *BotResourceRegistry) Add(msg *AddResourceMessage) {
	resource := Resource{ID: msg.ResourceID, Resource: msg.Resource, UserID: "unknown", PrototypeID: idGen(), Time: time.Now()}
	brr.botResources = append(brr.botResources, resource)
}

func (brr *BotResourceRegistry) Update(msg *UpdateResourceMessage) {
	for i, registryItem := range brr.botResources {
		if registryItem.ID == msg.ResourceID {
			brr.botResources[i] = Resource{ID: msg.ResourceID, Resource: msg.Resource}
		}
	}
}

func (brr *BotResourceRegistry) Remove(msg *RemoveResourceMessage) {
	for i, registryItem := range brr.botResources {
		if registryItem.ID == msg.ResourceID {
			brr.botResources = append(brr.botResources[:i], brr.botResources[i+1:]...)
		}
	}
}

func listen(wsConn *websocket.Conn) {
	defer wsConn.Close()

	for {
		_, b, err := wsConn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("NextReader:", err)
			}
			return
		}

		var f interface{}
		err = json.Unmarshal(b, &f)
		if err != nil {
			log.Printf("err: %s", err.Error())
		}

		msgMap := f.(map[string]interface{})
		if val, ok := msgMap["type"]; ok {
			if val == "auth" {
				wsConn.WriteJSON(NewAuthResultMessage(true, ""))

				// new connection, successfully auth'd, send all resources
				for _, botRes := range registry.botResources {
					fmt.Printf("%#v", botRes)
					newAddRes := NewAddResourceMessage(idGen(), &Resource{
						Resource: botRes.Resource,
					})
					wsConn.WriteJSON(newAddRes)
				}
			}
		} else {
			fmt.Printf("Invalid message received: %s\n", string(b))
		}
	}
}

// apiResourceHandler is used by the dev-console ui to upate the resource map, triggering
// a websocket message to be sent to the bot process.
func apiResourceHandler2(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}

	var postMsg ConsolePostMessage
	err = json.Unmarshal(body, &postMsg)
	if err != nil {
		log.Panicln("err: ", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case "POST":
		addResMsg := NewAddResourceMessage(postMsg.ResourceID, &postMsg.Resource)
		registry.Add(addResMsg)

		if wsConn != nil {
			wsConn.WriteJSON(addResMsg)
		}

		// respond to http call
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(""))

	case "PATCH":
		updateResMsg := NewUpdateResourceMessage(postMsg.ResourceID, &postMsg.Resource)
		registry.Update(updateResMsg)

		if wsConn != nil {
			wsConn.WriteJSON(updateResMsg)
		}

		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))

	case "DELETE":
		removeResMsg := NewRemoveResourceMessage(postMsg.ResourceID)
		registry.Remove(removeResMsg)

		if wsConn != nil {
			wsConn.WriteJSON(removeResMsg)
		}

		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))
	}

	fmt.Printf("registry:\n%s", registry)
}

func idGen() string {
	return uuid.NewRandom().String()
}
