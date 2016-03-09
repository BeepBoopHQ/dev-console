package main

import (
	"time"

	"github.com/pborman/uuid"
)

type Resource struct {
	ID          string `json:"resourceID"`
	UserID      string
	PrototypeID string
	ResourceCfg map[string]string `json:"resource"`
	Time        time.Time
}

// ============================================================================
// Message Types
// ============================================================================

// MessageMeta embedded type for all messages
type MessageMeta struct {
	Type string    `json:"type"`
	Date time.Time `json:"date"`
	ID   string    `json:"msgID"`
}

func NewMessageMeta(messageType string) *MessageMeta {
	return &MessageMeta{
		Type: messageType,
		Date: time.Now(),
		ID:   uuid.NewRandom().String(),
	}
}

// RemoveResourceMessage sent by server to notify client of an error
type ErrorMessage struct {
	MessageMeta
	Error string `json:"error"`
}

const ErrorMessageType = "error"

func NewErrorMessage(msg string) *ErrorMessage {
	return &ErrorMessage{
		MessageMeta: *NewMessageMeta(ErrorMessageType),
		Error:       msg,
	}
}

// AuthMessage send by client to initiate authentication
type AuthMessage struct {
	MessageMeta
	ID    string `json:"id"`
	Token string `json:"token"`
}

const AuthMessageType = "auth"

func NewAuthMessage(id, token string) *AuthMessage {
	return &AuthMessage{
		MessageMeta: *NewMessageMeta(AuthMessageType),
		ID:          id,
		Token:       token,
	}
}

//AuthResultMessage send by server to deliver authentication result
type AuthResultMessage struct {
	MessageMeta
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

const AuthResultType = "auth_result"

func NewAuthResultMessage(success bool, err string) *AuthResultMessage {
	return &AuthResultMessage{
		MessageMeta: *NewMessageMeta(AuthResultType),
		Success:     success,
		Error:       err,
	}
}

// AddResourceMessage send by server to notify client of a new resource to manage
type AddResourceMessage struct {
	MessageMeta
	ResourceID string            `json:"resourceID"`
	Resource   map[string]string `json:"resource"`
}

const AddResourceMessageType = "add_resource"

func NewAddResourceMessage(resourceID string, resource *Resource) *AddResourceMessage {
	return &AddResourceMessage{
		MessageMeta: *NewMessageMeta(AddResourceMessageType),
		ResourceID:  resourceID,
		Resource:    resource.ResourceCfg,
	}
}

// UpdateResourceMessage send by server to notify client of a new resource to manage
type UpdateResourceMessage struct {
	MessageMeta
	ResourceID string            `json:"resourceID"`
	Resource   map[string]string `json:"resource"`
}

const UpdateResourceMessageType = "update_resource"

func NewUpdateResourceMessage(resourceID string, resource *Resource) *UpdateResourceMessage {
	return &UpdateResourceMessage{
		MessageMeta: *NewMessageMeta(UpdateResourceMessageType),
		ResourceID:  resourceID,
		Resource:    resource.ResourceCfg,
	}
}

// RemoveResourceMessage send by server to notify client of a resource they should stop managing
type RemoveResourceMessage struct {
	MessageMeta
	ResourceID string `json:"resourceID"`
}

const RemoveResourceMessageType = "remove_resource"

func NewRemoveResourceMessage(resourceID string) *RemoveResourceMessage {
	return &RemoveResourceMessage{
		MessageMeta: *NewMessageMeta(RemoveResourceMessageType),
		ResourceID:  resourceID,
	}
}
