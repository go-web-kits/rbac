package rbac

import (
	"reflect"

	"github.com/pkg/errors"
)

// To Subject:
// TODO batch
func (a Assignment) notBeA(temporary bool, role interface{}, args ...interface{}) error {
	config := ConfigOf(a.Object)
	foundSubject, err := config.Engine.FindSubject(a.Object, args...)
	if err != nil {
		return err // errors.Wrapf(err, "rbac.Let->FindSubject") TODO
	}
	tempSubject, isTempSubject := foundSubject.(TemporaryRoleAssignable)

	roles := []interface{}{}
	refRole := reflect.Indirect(reflect.ValueOf(role))
	if refRole.Kind() != reflect.Slice {
		refRole = reflect.Indirect(reflect.ValueOf([]interface{}{role}))
	}
	for i := 0; i < refRole.Len(); i++ {
		foundRole, err := config.Engine.FindRole(refRole.Index(i).Interface(), args...)
		if err != nil {
			return err // errors.Wrapf(err, "rbac.Let->FindRole")
		}
		if !Whether(foundSubject).IsA(foundRole) {
			return errors.Wrapf(NotAssigment, "rbac.Let->NotBeA Role: %#v", foundRole)
		}

		if temporary && isTempSubject && tempSubject.IsTemporaryRole(foundRole) {
			if err := tempSubject.TemporaryRoleRemove(foundRole); err != nil {
				return err
			}
		}
		roles = append(roles, foundRole)
	}

	if !temporary {
		if err := config.Engine.CancelAssignRole(foundSubject, roles, args...); err != nil {
			return err
		}
		if callbackable, ok := foundSubject.(AssignCallbacks); ok {
			callbackable.AfterCancelAssign(roles, args...)
		}
	}

	return nil
}

func (a Assignment) NotBeA(role Role, args ...interface{}) error {
	return a.notBeA(false, role, args...)
}
func (a Assignment) NotBecomeA(role Role, args ...interface{}) error {
	return a.notBeA(false, role, args...)
}
func (a Assignment) NotBeATemp(role Role, args ...interface{}) error {
	return a.notBeA(true, role, args...)
}

// To Role:
// TODO batch
func (a Assignment) Cannot(permission Permission, args ...interface{}) error {
	config := ConfigOf(a.Object)
	foundRole, err := config.Engine.FindRole(a.Object, args...)
	if err != nil {
		return err // errors.Wrapf(err, "rbac.Let->FindRole")
	}

	permissions := []interface{}{}
	refPms := reflect.Indirect(reflect.ValueOf(permission))
	if refPms.Kind() != reflect.Slice {
		refPms = reflect.Indirect(reflect.ValueOf([]interface{}{permission}))
	}
	for i := 0; i < refPms.Len(); i++ {
		foundPms, err := config.Engine.FindPermission(refPms.Index(i).Interface(), args...)
		if err != nil {
			return err // errors.Wrapf(err, "rbac.Let->FindPermission")
		}
		if !Whether(foundRole).Can(foundPms) {
			return errors.Wrapf(NotAssigment, "rbac.Let->Cannot Permission: %#v", foundPms)
		}
		permissions = append(permissions, foundPms)
	}

	if err := config.Engine.CancelAssignPermission(foundRole, permissions, args...); err != nil {
		return err
	}

	if callbackable, ok := foundRole.(AssignCallbacks); ok {
		callbackable.AfterCancelAssign(permissions, args...)
	}
	return nil
}

// To Resource:
// TODO batch
func (a Assignment) CannotBe(action string) Assignment {
	resource, ok := a.Object.(Resource)
	if !ok {
		panic("rbac.Let->CanBe: Object is not a resource")
	}

	return Assignment{permissions: []Permission{resource.NewPermissionForIt(action)}, cancel: true}
}
