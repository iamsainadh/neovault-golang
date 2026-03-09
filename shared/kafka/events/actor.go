package events

type ActorType string

const (
	ActorUser           ActorType = "USER"
	ActorSystem         ActorType = "SYSTEM"
	ActorServiceAccount ActorType = "SERVICE_ACCOUNT"
)

type Actor struct {
	Type ActorType `json:"type"`
	ID   string    `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}
