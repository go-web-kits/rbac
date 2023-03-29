package rbac_dbx

import (
	"github.com/go-web-kits/dbx"
	"github.com/go-web-kits/rbac"
)

func (e Engine) AssignRole(foundSubject interface{}, foundRoles []interface{}, opts ...interface{}) error {
	assocName := rbac.ConfigOf(foundSubject).Roles()
	return dbx.ScopingModel(foundSubject, getOpts(opts...)...).Association(assocName).Append(foundRoles...).Error
}

func (e Engine) AssignPermission(foundRole interface{}, foundPms []interface{}, opts ...interface{}) error {
	assocName := rbac.ConfigOf(foundRole).Permissions()
	return dbx.ScopingModel(foundRole, getOpts(opts...)...).Association(assocName).Append(foundPms...).Error
}

func (e Engine) CancelAssignRole(foundSubject interface{}, foundRole []interface{}, opts ...interface{}) error {
	assocName := rbac.ConfigOf(foundSubject).Roles()
	return dbx.ScopingModel(foundSubject, getOpts(opts...)...).Association(assocName).Delete(foundRole...).Error
}

func (e Engine) CancelAssignPermission(foundRole interface{}, foundPms []interface{}, opts ...interface{}) error {
	assocName := rbac.ConfigOf(foundRole).Permissions()
	return dbx.ScopingModel(foundRole, getOpts(opts...)...).Association(assocName).Delete(foundPms...).Error
}

func (e Engine) ClearRole(foundSubject interface{}, opts ...interface{}) error {
	assocName := rbac.ConfigOf(foundSubject).Roles()
	return dbx.ScopingModel(foundSubject, getOpts(opts...)...).Association(assocName).Clear().Error
}

func (e Engine) ClearPermission(foundRole interface{}, opts ...interface{}) error {
	assocName := rbac.ConfigOf(foundRole).Permissions()
	return dbx.ScopingModel(foundRole, getOpts(opts...)...).Association(assocName).Clear().Error
}
