package test

import (
	. "github.com/go-web-kits/rbac"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("RbacDBx", func() {
	DefinitionTesting()
	AssignmentTesting()
	QueryingTesting()
})

type People struct {
	ID   uint   `json:"id" db:"id" gorm:"primary_key;index"`
	Name string `json:"name" db:"name"`

	PeopleRoles []PeopleRole `gorm:"many2many:peoples_and_people_roles"`
	TempRoles   []PeopleRole
	ActAsSubject
}

type PeopleRole struct {
	ID   uint   `json:"id" db:"id" gorm:"primary_key;index"`
	Name string `json:"name" db:"name"`

	PeoplePermissions []PeoplePermission `gorm:"many2many:people_roles_and_people_permissions"`
	ActAsRole
}

type PeoplePermission struct {
	ID       uint   `json:"id" db:"id" gorm:"primary_key;index"`
	Action   string `json:"action" db:"action"`
	RecordID uint   `json:"record_id" db:"record_id"`

	ActAsPermission
}

type PeopleResource struct {
	ID uint `json:"id" db:"id" gorm:"primary_key;index"`

	ActAsResource
}

func (p *People) TemporaryRoleAssign(foundRole interface{}) error {
	p.TempRoles = append(p.TempRoles, foundRole.(PeopleRole))
	return nil
}

func (p *People) TemporaryRoleRemove(foundRole interface{}) error {
	tRoles := []PeopleRole{}
	for _, r := range p.TempRoles {
		if r.ID != foundRole.(PeopleRole).ID {
			tRoles = append(tRoles, r)
		}
	}

	p.TempRoles = tRoles
	return nil
}

func (p *People) TemporaryRoleClear() error {
	p.TempRoles = []PeopleRole{}
	return nil
}

func (p People) GetTemporaryRoles() (list []interface{}) {
	for _, r := range p.TempRoles {
		list = append(list, r)
	}
	return list
}

func (p People) IsTemporaryRole(foundRole interface{}) bool {
	for _, r := range p.TempRoles {
		if r.ID == foundRole.(PeopleRole).ID {
			return true
		}
	}
	return false
}

func (r PeopleResource) NewPermissionForIt(action string, args ...interface{}) Permission {
	return &PeoplePermission{Action: action, RecordID: r.ID}
}
