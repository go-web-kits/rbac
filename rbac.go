package rbac

// Subject
type Subject interface{ IsSubject() }
type TemporaryRoleAssignable interface {
	TemporaryRoleAssign(foundRole interface{}) error
	TemporaryRoleRemove(foundRole interface{}) error
	TemporaryRoleClear() error
	TemporaryRoleQueryable
}
type TemporaryRoleQueryable interface {
	GetTemporaryRoles() []interface{}
	IsTemporaryRole(foundRole interface{}) bool
}
type ActAsSubject struct{}

// Role
type Role interface{ IsRole() }
type ActAsRole struct{}

// Permission
type Permission interface{ IsPermission() }
type ActAsPermission struct{}

// Resource
type Resource interface {
	IsResource()
	NewPermissionForIt(action string, args ...interface{}) Permission
}
type ActAsResource struct{}

func (s ActAsSubject) IsSubject()       {}
func (r ActAsRole) IsRole()             {}
func (p ActAsPermission) IsPermission() {}
func (r ActAsResource) IsResource()     {}

type AssignCallbacks interface {
	AfterAssign(objs []interface{}, args ...interface{})
	AfterCancelAssign(objs []interface{}, args ...interface{})
	AfterClearAssign(args ...interface{})
}
