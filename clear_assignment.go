package rbac

func (a Assignment) NotBeAnyone(args ...interface{}) error {
	config := ConfigOf(a.Object)
	foundSubject, err := config.Engine.FindSubject(a.Object)
	if err != nil {
		return err // errors.Wrapf(err, "rbac.Let->FindSubject") TODO
	}

	err = config.Engine.ClearRole(foundSubject, args...)
	if err != nil {
		return err
	}

	if s, ok := foundSubject.(TemporaryRoleAssignable); ok {
		if err := s.TemporaryRoleClear(); err != nil {
			return err
		}
	}

	if callbackable, ok := foundSubject.(AssignCallbacks); ok {
		callbackable.AfterClearAssign(args...)
	}
	return nil
}

func (a Assignment) CannotDoAnything(args ...interface{}) error {
	config := ConfigOf(a.Object)
	foundRole, err := config.Engine.FindRole(a.Object)
	if err != nil {
		return err // errors.Wrapf(err, "rbac.Let->FindRole")
	}

	if err := config.Engine.ClearPermission(foundRole, args...); err != nil {
		return err
	}

	if callbackable, ok := foundRole.(AssignCallbacks); ok {
		callbackable.AfterClearAssign(args...)
	}
	return nil
}

func (a Assignment) ByNobody(args ...interface{}) error {
	if a.permissions == nil {
		panic("rbac.Let->By: wrong method calling, make sure you call `CanBe` before")
	}

	// return Assignment{Object: role}.Cannot(a.permissions, args...)
	panic("WIP")
}
