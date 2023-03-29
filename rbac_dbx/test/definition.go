package test

import (
	"github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/rbac"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var DefinitionTesting = func() {
	Describe("Definition", func() {
		var (
			tom       People
			admin     PeopleRole
			read_data PeoplePermission
		)

		BeforeEach(func() {
			tom = People{Name: "Tom"}
			admin = PeopleRole{Name: "admin"}
			read_data = PeoplePermission{Action: "read"}
		})

		AfterEach(func() {
			factory.Association(&tom, "PeopleRoles").Clear()
			factory.Association(&admin, "PeoplePermissions").Clear()
			CleanData(&People{}, &PeopleRole{}, &PeoplePermission{})
		})

		It("creates role for subject", func() {
			Expect(dbx.Model(&PeopleRole{}).Count()).To(BeZero())
			Expect(For(&People{}).Create(&admin)).To(Succeed())
			Expect(dbx.Model(&PeopleRole{}).Count()).To(BeEquivalentTo(1))
			Expect(dbx.Find(&admin, admin)).To(HaveFound())
		})

		It("creates permission for role", func() {
			Expect(dbx.Model(&PeoplePermission{}).Count()).To(BeZero())
			Expect(For(&PeopleRole{}).Create(&read_data)).To(Succeed())
			Expect(dbx.Model(&PeoplePermission{}).Count()).To(BeEquivalentTo(1))
			Expect(dbx.Find(&read_data, read_data)).To(HaveFound())
		})
	})
}
