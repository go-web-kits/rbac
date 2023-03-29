package rbac

import "github.com/go-web-kits/utils"

type Definition struct {
	Subject Subject
	Role    Role
}

func For(obj interface{}) Definition {
	if subject, ok := obj.(Subject); ok {
		return Definition{Subject: subject}
	} else if role, ok := obj.(Role); ok {
		return Definition{Role: role}
	} else {
		panic("rbac.For: object `" + utils.TypeNameOf(obj) + "` is not a Subject or Role")
	}
}

// TODO batch
func (d Definition) Create(obj interface{}, opts ...interface{}) error {
	if role, ok := obj.(Role); ok {
		if d.Subject == nil {
			panic("rbac.For->CreateRole: object `" + utils.TypeNameOf(d.Role) + "` is not a Subject")
		}

		config := ConfigOf(d.Subject)
		if _, err := config.Engine.FindRole(role); err == nil {
			return nil
		}
		return config.Engine.CreateRole(d.Subject, role, opts...)
	} else if permission, ok := obj.(Permission); ok {
		if d.Role == nil {
			panic("rbac.For->CreatePermission: object `" + utils.TypeNameOf(d.Subject) + "` is not a Role")
		}

		config := ConfigOf(d.Role)
		if _, err := config.Engine.FindPermission(permission); err == nil {
			return nil
		}
		return EngineOf(d.Role).CreatePermission(d.Role, permission, opts...)
	} else {
		panic("rbac.For->Create: object `" + utils.TypeNameOf(obj) + "` is not a Role or Permission")
	}
}

func (d Definition) CreateRole(role Role, opts ...interface{}) error {
	return d.Create(role, opts...)
}

func (d Definition) CreatePermission(permission Permission, opts ...interface{}) error {
	return d.Create(permission, opts...)
}
