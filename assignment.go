package rbac

import (
	"reflect"

	"github.com/pkg/errors"
)

type Assignment struct {
	Object      interface{}
	permissions []Permission
	cancel      bool
	only        bool
}

func Let(obj interface{}) Assignment {
	return Assignment{Object: obj}
}

// To Subject:
func (a Assignment) beA(temporary bool, role interface{}, args ...interface{}) error {
	config := ConfigOf(a.Object)
	foundSubject, err := config.Engine.FindSubject(a.Object, args...)
	if err != nil {
		return err // errors.Wrapf(err, "rbac.Let->FindSubject")
	}
	tempSubject, isTempSubject := foundSubject.(TemporaryRoleAssignable)

	roles := []interface{}{}
	refRole := reflect.Indirect(reflect.ValueOf(role))
	if refRole.Kind() != reflect.Slice {
		refRole = reflect.Indirect(reflect.ValueOf([]interface{}{role}))
	}
	for i := 0; i < refRole.Len(); i++ {
		foundRole, err := findOrCreateRole(refRole.Index(i).Interface(), config, args...)
		if err != nil {
			return err
		}
		if Whether(foundSubject).IsA(foundRole) {
			continue
		}

		if temporary && isTempSubject && !tempSubject.IsTemporaryRole(foundRole) {
			if err := tempSubject.TemporaryRoleAssign(foundRole); err != nil {
				return err
			}
		}
		roles = append(roles, foundRole)
	}

	if !temporary {
		if err := config.Engine.AssignRole(foundSubject, roles, args...); err != nil {
			return err
		}
		if callbackable, ok := foundSubject.(AssignCallbacks); ok {
			callbackable.AfterAssign(roles, args...)
		}
	}

	return nil
}

func (a Assignment) BeA(role interface{}, args ...interface{}) error {
	return a.beA(false, role, args...)
}
func (a Assignment) BecomeA(role interface{}, args ...interface{}) error {
	return a.beA(false, role, args...)
}
func (a Assignment) BeATemp(role interface{}, args ...interface{}) error {
	return a.beA(true, role, args...)
}

func (a Assignment) OnlyBeA(role interface{}, args ...interface{}) error {
	if err := a.NotBeAnyone(args...); err != nil {
		return err
	}
	return a.BeA(role, args...)
}
func (a Assignment) OnlyBecomeA(role interface{}, args ...interface{}) error {
	return a.OnlyBeA(role, args...)
}

// To Role:
func (a Assignment) Can(permission interface{}, args ...interface{}) error {
	config := ConfigOf(a.Object)
	foundRole, err := findOrCreateRole(a.Object, config, args...)
	if err != nil {
		return err
	}

	permissions := []interface{}{}
	refPms := reflect.Indirect(reflect.ValueOf(permission))
	if refPms.Kind() != reflect.Slice {
		refPms = reflect.Indirect(reflect.ValueOf([]interface{}{permission}))
	}
	for i := 0; i < refPms.Len(); i++ {
		foundPms, err := findOrCreatePermission(refPms.Index(i).Interface(), config, args...)
		if err != nil {
			return err
		}
		if Whether(foundRole).Can(foundPms) {
			continue
		}
		permissions = append(permissions, foundPms)
	}

	if err := config.Engine.AssignPermission(foundRole, permissions, args...); err != nil {
		return err
	}

	if callbackable, ok := foundRole.(AssignCallbacks); ok {
		callbackable.AfterAssign(permissions, args...)
	}
	return nil
}

func (a Assignment) CanOnly(permission interface{}, args ...interface{}) error {
	if err := a.CannotDoAnything(args...); err != nil {
		return err
	}
	return a.Can(permission, args...)
}

// To Resource:
func (a Assignment) CanBe(actions ...string) Assignment {
	resource, ok := a.Object.(Resource)
	if !ok {
		panic("rbac.Let->CanBe: Object is not a resource")
	}

	permissions := []Permission{}
	for _, action := range actions {
		permissions = append(permissions, resource.NewPermissionForIt(action))
	}
	return Assignment{permissions: permissions}
}

func (a Assignment) CanOnlyBe(actions ...string) Assignment {
	assignment := a.CanBe(actions...)
	assignment.only = true
	return assignment
}

func (a Assignment) By(role Role, args ...interface{}) error {
	if a.permissions == nil {
		panic("rbac.Let->By: wrong method calling, make sure you call `CanBe` before")
	}

	if a.cancel {
		return Assignment{Object: role}.Cannot(a.permissions[0], args...)
	} else if a.only {
		return Assignment{Object: role}.CanOnly(a.permissions, args...)
	} else {
		return Assignment{Object: role}.Can(a.permissions, args...)
	}
}

// ===============================

func findOrCreateRole(role interface{}, config Config, args ...interface{}) (Role, error) {
	foundRole, err := config.Engine.FindRole(role, args...)
	if err != nil && config.AutoDefineRole {
		if err := config.Engine.CreateRole(nil, foundRole, args...); err != nil {
			return nil, errors.Wrapf(err, "rbac.Let->AutoDefineRole")
		}
	} else if err != nil {
		return nil, err // errors.Wrapf(err, "rbac.Let->FindRole")
	}
	return foundRole, nil
}

func findOrCreatePermission(permission interface{}, config Config, args ...interface{}) (Permission, error) {
	foundPms, err := config.Engine.FindPermission(permission, args...)
	if err != nil && config.AutoDefinePermission {
		if err := config.Engine.CreatePermission(nil, foundPms, args...); err != nil {
			return nil, errors.Wrapf(err, "rbac.Let->AutoDefinePermission")
		}
	} else if err != nil {
		return nil, err // errors.Wrapf(err, "rbac.Let->FindPermission")
	}
	return foundPms, nil
}
