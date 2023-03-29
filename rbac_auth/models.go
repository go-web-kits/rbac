package rbac_auth

import "github.com/go-web-kits/rbac"

type AuthSubject struct {
	AD    string `json:"ad"`
	Token string `json:"token"`

	Roles []Role
	rbac.ActAsSubject
}

type Subject = AuthSubject

func (s Subject) AllowedResourceNames(action string) []string {
	resources := []string{}
	for _, role := range s.Roles {
		for _, permission := range role.Permissions {
			if permission.Action == action {
				resources = append(resources, permission.Resource.Name)
			}
		}
	}
	return resources
}

type AuthRole struct {
	Name string `json:"name"`
	Desc string `json:"description"`

	Permissions []Permission `json:"abilities"`
	rbac.ActAsRole
}

type Role = AuthRole

type AuthPermission struct {
	Action string `json:"abilityCode"`
	Desc   string `json:"description"`

	Resource Resource
	rbac.ActAsPermission
}

type Permission = AuthPermission

type AuthResource struct {
	Name string `json:"name"`

	rbac.ActAsResource
}

type Resource = AuthResource

func (r AuthResource) NewPermissionForIt(action string, args ...interface{}) rbac.Permission {
	return AuthPermission{Action: action, Resource: r}
}
