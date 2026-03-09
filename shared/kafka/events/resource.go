package events

type ResourceType string

const (
	ResourceSecret       ResourceType = "SECRET"
	ResourceVault        ResourceType = "VAULT"
	ResourceOrganization ResourceType = "ORGANIZATION"
	ResourceWorkspace    ResourceType = "WORKSPACE"
	ResourceProject      ResourceType = "PROJECT"
	ResourceUser         ResourceType = "USER"
	ResourceAPIKey       ResourceType = "API_KEY"
)

type Resource struct {
	Type ResourceType `json:"type"`
	ID   string       `json:"id,omitempty"`
	Name string       `json:"name,omitempty"`
}
