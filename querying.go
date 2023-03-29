package rbac

type Querying struct {
	Object     interface{}
	permission Permission
}

func Whether(obj interface{}) Querying {
	return Querying{Object: obj}
}

// TODO: batch
func (q Querying) IsA(role interface{}, roles ...interface{}) bool {
	config := ConfigOf(q.Object)
	foundSubject, err := config.Engine.FindSubject(q.Object)
	if err != nil {
		return false
	}
	foundRole, err := config.Engine.FindRole(role)
	if err != nil {
		return false
	}

	if s, ok := q.Object.(TemporaryRoleQueryable); ok {
		return s.IsTemporaryRole(foundRole) || config.Engine.QueryIsRole(foundSubject, foundRole)
	} else {
		return config.Engine.QueryIsRole(foundSubject, foundRole)
	}
}

// TODO
func (q Querying) IsOneOfThe(role interface{}, roles ...interface{}) bool {
	panic("WIP")
}

// TODO: batch
func (q Querying) Can(permission interface{}, permissions ...interface{}) bool {
	config := ConfigOf(q.Object)
	var foundObj interface{}
	var err error
	if s, ok := q.Object.(Subject); ok {
		foundObj, err = config.Engine.FindSubject(s) // TODO Find or Reload?
		if err != nil {
			return false
		}
	} else if r, ok := q.Object.(Role); ok {
		foundObj, err = config.Engine.FindRole(r)
		if err != nil {
			return false
		}
	} else {
		panic("rbac.Whether->Can: Object is not Role or Subject")
	}

	foundPms, err := config.Engine.FindPermission(permission)
	if err != nil {
		return false
	}

	if s, ok := q.Object.(TemporaryRoleQueryable); ok {
		for _, role := range s.GetTemporaryRoles() {
			if config.Engine.QueryHavePermission(role, foundPms) {
				return true
			}
		}
	}
	return config.Engine.QueryHavePermission(foundObj, foundPms)
}

// TODO
func (q Querying) CanOneOf(permission interface{}, permissions ...interface{}) bool {
	panic("WIP")
}

// TODO batch
func (q Querying) CanBe(action string, actions ...string) Querying {
	resource, ok := q.Object.(Resource)
	if !ok {
		panic("rbac.Whether->CanBe: Object is not a resource")
	}

	return Querying{permission: resource.NewPermissionForIt(action)}
}

// TODO
func (q Querying) CanBeOneOf(action string, actions ...string) Querying {
	panic("WIP")
}

func (q Querying) By(obj interface{}) bool {
	if q.permission == nil {
		panic("rbac.Whether->By: wrong method calling, make sure you call `CanBe` before")
	}

	return Querying{Object: obj}.Can(q.permission)
}
