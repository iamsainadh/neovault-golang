package events

type Action string

const (
	ActionCreateSecret Action = "CREATE_SECRET"
	ActionUpdateSecret Action = "UPDATE_SECRET"
	ActionDeleteSecret Action = "DELETE_SECRET"
	ActionReadSecret   Action = "READ_SECRET"

	ActionCreateVault Action = "CREATE_VAULT"
	ActionUpdateVault Action = "UPDATE_VAULT"
	ActionDeleteVault Action = "DELETE_VAULT"

	ActionLogin  Action = "LOGIN"
	ActionLogout Action = "LOGOUT"

	ActionCreateOrganization Action = "CREATE_ORGANIZATION"
	ActionUpdateOrganization Action = "UPDATE_ORGANIZATION"
	ActionDeleteOrganization Action = "DELETE_ORGANIZATION"

	ActionInviteUser Action = "INVITE_USER"
	ActionRemoveUser Action = "REMOVE_USER"

	ActionRotateAPIKey Action = "ROTATE_API_KEY"
	ActionAccessDenied Action = "ACCESS_DENIED"
)
