package rbac_dbx

import (
	"reflect"

	"github.com/go-web-kits/dbx"
	"github.com/go-web-kits/rbac"
)

func (e Engine) QueryIsRole(foundSubject interface{}, foundRole interface{}) bool {
	config := rbac.ConfigOf(foundSubject)
	roles := reflect.Indirect(reflect.ValueOf(foundSubject)).FieldByName(config.Roles())
	return include(foundRole, roles)
}

func (e Engine) QueryHavePermission(obj interface{}, foundPms interface{}) bool {
	config := rbac.ConfigOf(obj)
	refObj := reflect.Indirect(reflect.ValueOf(obj))

	if _, isSubject := obj.(rbac.Subject); isSubject {
		roles := refObj.FieldByName(config.Roles())
		for i := 0; i < roles.Len(); i++ {
			if include(foundPms, roles.Index(i).FieldByName(config.Permissions())) {
				return true
			}
		}
		return false
	} else if _, isRole := obj.(rbac.Role); isRole {
		return include(foundPms, refObj.FieldByName(config.Permissions()))
	} else {
		panic("rbac.QueryHavePermission: Object is not Subject or Role")
	}
}

// =======

func include(obj interface{}, objs reflect.Value) bool {
	id := dbx.IdOf(obj)
	for i := 0; i < objs.Len(); i++ {
		if id == dbx.IdOf(objs.Index(i).Interface()) {
			return true
		}
	}
	return false
}
