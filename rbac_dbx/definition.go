package rbac_dbx

import "github.com/go-web-kits/dbx"

func (e Engine) CreateRole(subject interface{}, role interface{}, opts ...interface{}) error {
	return dbx.Create(role, getOpts(opts...)...).Err
}

func (e Engine) CreatePermission(role interface{}, permission interface{}, opts ...interface{}) error {
	return dbx.Create(permission, getOpts(opts...)...).Err
}
