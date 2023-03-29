package rbac_dbx

import (
	"github.com/go-web-kits/dbx"
	"github.com/go-web-kits/rbac"
	"github.com/go-web-kits/utils/structx"
)

func (e Engine) FindSubject(obj interface{}, opts ...interface{}) (found rbac.Subject, err error) {
	subject := obj.(rbac.Subject) // TODO: support string
	if !(dbx.Result{Data: subject}).IsNewRecord() {
		return subject, nil
	}

	config := rbac.ConfigOf(subject)
	preload := config.Roles() + "." + config.Permissions()
	err = dbx.Find(subject, structx.From(subject).ToMap().Compact(), dbx.With{Preload: preload}).Err
	if err != nil {
		return subject, err
	}
	return subject, nil
}

func (e Engine) FindRole(obj interface{}, opts ...interface{}) (found rbac.Role, err error) {
	role := obj.(rbac.Role) // TODO: support string
	if !(dbx.Result{Data: role}).IsNewRecord() {
		return role, nil
	}

	err = dbx.Find(role, structx.From(role).ToMap().Compact(), dbx.With{Preload: rbac.ConfigOf(role).Permissions()}).Err
	if err != nil {
		return role, err
	}
	return role, nil
}

func (e Engine) FindPermission(obj interface{}, opts ...interface{}) (found rbac.Permission, err error) {
	pms := obj.(rbac.Permission)
	if !(dbx.Result{Data: pms}).IsNewRecord() {
		return pms, nil
	}
	if err = dbx.Find(pms, structx.From(pms).ToMap().Compact()).Err; err != nil {
		return pms, err
	}
	return pms, nil
}
