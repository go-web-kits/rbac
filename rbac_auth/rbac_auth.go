package rbac_auth

import "github.com/go-web-kits/rbac"

type Engine struct{}

func (e Engine) QueryIsRole(foundSubject interface{}, foundRole interface{}) bool {
	subject, role := foundSubject.(Subject), foundRole.(Role)
	for _, r := range subject.Roles {
		if role.Name == r.Name {
			return true
		}
	}
	return false
}

func (e Engine) QueryHavePermission(obj interface{}, foundPms interface{}) bool {
	permission, permissions := foundPms.(Permission), []Permission{}
	if subject, ok := obj.(Subject); ok {
		for _, role := range subject.Roles {
			permissions = append(permissions, role.Permissions...)
		}
	} else if role, ok := obj.(Role); ok {
		permissions = role.Permissions
	}

	for _, p := range permissions {
		if p.Action == permission.Action && p.Resource == permission.Resource {
			return true
		}
	}
	return false
}

func (e Engine) FindSubject(obj interface{}, args ...interface{}) (found rbac.Subject, err error) {
	return obj.(Subject), nil // TODO
}

func (e Engine) FindRole(obj interface{}, args ...interface{}) (found rbac.Role, err error) {
	if s, ok := obj.(string); ok {
		return Role{Name: s}, nil
	}
	return obj.(Role), nil // TODO
}

func (e Engine) FindPermission(obj interface{}, args ...interface{}) (found rbac.Permission, err error) {
	return obj.(Permission), nil // TODO
}

func (e Engine) CreateRole(subject interface{}, role interface{}, args ...interface{}) error {
	panic("WIP")
}

func (e Engine) CreatePermission(role interface{}, permission interface{}, args ...interface{}) error {
	panic("WIP")
}

func (e Engine) AssignRole(foundSubject interface{}, foundRoles []interface{}, args ...interface{}) error {
	panic("WIP")
}

func (e Engine) AssignPermission(foundRole interface{}, foundPms []interface{}, args ...interface{}) error {
	panic("WIP")
}

func (e Engine) CancelAssignRole(foundSubject interface{}, foundRole []interface{}, args ...interface{}) error {
	panic("WIP")
}

func (e Engine) CancelAssignPermission(foundRole interface{}, foundPms []interface{}, args ...interface{}) error {
	panic("WIP")
}

func (e Engine) ClearRole(foundSubject interface{}, args ...interface{}) error {
	panic("WIP")
}

func (e Engine) ClearPermission(foundRole interface{}, args ...interface{}) error {
	panic("WIP")
}
