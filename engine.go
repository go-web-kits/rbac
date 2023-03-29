package rbac

type Engine interface {
	CreateRole(subject interface{}, role interface{}, args ...interface{}) error
	CreatePermission(role interface{}, permission interface{}, args ...interface{}) error

	AssignRole(foundSubject interface{}, foundRoles []interface{}, args ...interface{}) error
	AssignPermission(foundRole interface{}, foundPms []interface{}, args ...interface{}) error

	CancelAssignRole(foundSubject interface{}, foundRoles []interface{}, args ...interface{}) error
	CancelAssignPermission(foundRole interface{}, foundPms []interface{}, args ...interface{}) error
	ClearRole(foundSubject interface{}, args ...interface{}) error
	ClearPermission(foundRole interface{}, args ...interface{}) error

	QueryIsRole(foundSubject interface{}, foundRole interface{}) bool
	QueryHavePermission(obj interface{}, foundPms interface{}) bool

	FindSubject(obj interface{}, args ...interface{}) (found Subject, err error)
	// ReloadSubject(obj interface{}) (found Subject, err error)
	FindRole(obj interface{}, args ...interface{}) (found Role, err error)
	FindPermission(obj interface{}, args ...interface{}) (found Permission, err error)
}

func EngineOf(obj interface{}) (engine Engine) {
	config := ConfigOf(obj)
	engine = config.Engine
	if engine == nil {
		panic("Cannot Get Rbac Engine, Please Check Rbac Configuration!")
	}
	return engine
}
