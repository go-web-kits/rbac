package test

import (
	"github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/rbac"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var AssignmentTesting = func() {
	var (
		tom       People
		admin     PeopleRole
		read_data PeoplePermission
		resource  PeopleResource
	)

	BeforeEach(func() {
		tom = People{Name: "Tom"}
		admin = PeopleRole{Name: "admin"}
		read_data = PeoplePermission{Action: "read"}
		resource = PeopleResource{}
	})

	AfterEach(func() {
		factory.Association(&tom, "PeopleRoles").Clear()
		factory.Association(&admin, "PeoplePermissions").Clear()
		CleanData(&People{}, &PeopleRole{}, &PeoplePermission{})
	})

	Describe("Assignment", func() {
		BeforeEach(func() {
			factory.Create(&tom)
			_ = For(&People{}).Create(&admin)
			_ = For(&PeopleRole{}).Create(&read_data)
		})

		When("Role", func() {
			BeforeEach(func() {
				factory.Reload(&tom, dbx.With{Preload: "PeopleRoles"})
				Expect(tom.PeopleRoles).To(BeEmpty())
			})

			It("can assign the role to the subject successfully", func() {
				Expect(Let(&tom).BeA(&admin)).To(Succeed())
				factory.Reload(&tom, dbx.With{Preload: "PeopleRoles"})
				Expect(tom.PeopleRoles).To(Equal([]PeopleRole{admin}))
			})

			It("can cancel assign the role to the subject successfully", func() {
				Expect(Let(&tom).BeA(&admin)).To(Succeed())
				Expect(Let(&tom).NotBeA(&admin)).To(Succeed())
				factory.Reload(&tom, dbx.With{Preload: "PeopleRoles"})
				Expect(tom.PeopleRoles).To(BeEmpty())

				Expect(Let(&tom).BeA(&admin)).To(Succeed())
				Expect(Let(&tom).NotBeAnyone()).To(Succeed())
				factory.Reload(&tom, dbx.With{Preload: "PeopleRoles"})
				Expect(tom.PeopleRoles).To(BeEmpty())
			})

			It("should auto load", func() {
				Expect(Let(&People{Name: "Tom"}).BeA(&PeopleRole{Name: "admin"})).To(Succeed())
				factory.Reload(&tom, dbx.With{Preload: "PeopleRoles"})
				Expect(tom.PeopleRoles).To(Equal([]PeopleRole{admin}))
			})

			When("the subject is a TemporaryRoleAssignable", func() {
				It("can temporary assign the role to the subject successfully", func() {
					Expect(tom.IsTemporaryRole(admin)).To(BeFalse())
					Expect(Let(&tom).BeATemp(admin)).To(Succeed())
					Expect(tom.IsTemporaryRole(admin)).To(BeTrue())
					Expect(Let(&tom).NotBeATemp(admin)).To(Succeed())
					Expect(tom.IsTemporaryRole(admin)).To(BeFalse())
				})
			})
		})

		When("Permission", func() {
			BeforeEach(func() {
				factory.Reload(&admin, dbx.With{Preload: "PeoplePermissions"})
				Expect(admin.PeoplePermissions).To(BeEmpty())
			})

			It("can assign the permission to the role successfully", func() {
				Expect(Let(&admin).Can(&read_data)).To(Succeed())
				factory.Reload(&admin, dbx.With{Preload: "PeoplePermissions"})
				Expect(admin.PeoplePermissions).To(Equal([]PeoplePermission{read_data}))
			})

			It("can cancel assign the permission to the role successfully", func() {
				Expect(Let(&admin).Can(&read_data)).To(Succeed())
				Expect(Let(&admin).Cannot(&read_data)).To(Succeed())
				factory.Reload(&admin, dbx.With{Preload: "PeoplePermissions"})
				Expect(admin.PeoplePermissions).To(BeEmpty())

				Expect(Let(&admin).Can(&read_data)).To(Succeed())
				Expect(Let(&admin).CannotDoAnything()).To(Succeed())
				factory.Reload(&admin, dbx.With{Preload: "PeoplePermissions"})
				Expect(admin.PeoplePermissions).To(BeEmpty())
			})

			It("should auto load", func() {
				Expect(Let(&PeopleRole{Name: "admin"}).Can(&PeoplePermission{Action: "read"})).To(Succeed())
				factory.Reload(&admin, dbx.With{Preload: "PeoplePermissions"})
				Expect(admin.PeoplePermissions).To(Equal([]PeoplePermission{read_data}))
			})

			Context("Resource", func() {
				BeforeEach(func() {
					factory.Create(&resource)
					factory.Create(resource.NewPermissionForIt("read"))
				})

				It("can new the permission and assign it to the role", func() {
					Expect(Let(&resource).CanBe("read").By(&admin)).To(Succeed())
					factory.Reload(&admin, dbx.With{Preload: "PeoplePermissions"})
					Expect(admin.PeoplePermissions[0]).To(HaveAttributes(PeoplePermission{Action: "read", RecordID: resource.ID}))
				})
			})
		})
	})
}
