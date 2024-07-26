package users

type UserRole string

const (
	RoleMaster UserRole = "MASTER"
	RoleAdmin  UserRole = "ADMIN"
	RoleUser   UserRole = "USER"
)

type PermissionAction string

const (
	ActionCreate PermissionAction = "CREATE"
	ActionRead   PermissionAction = "READ"
	ActionUpdate PermissionAction = "UPDATE"
	ActionDelete PermissionAction = "DELETE"
)


