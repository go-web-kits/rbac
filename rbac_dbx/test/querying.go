package test

import (
	. "github.com/go-web-kits/rbac"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var QueryingTesting = func() {
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

	Describe("Querying", func() {
		BeforeEach(func() {
			factory.Create(&tom)
			_ = For(&People{}).Create(&admin)
			_ = For(&PeopleRole{}).Create(&read_data)
		})

		Describe("Role", func() {
			It("returns the assignment status correctly", func() {
				Expect(Whether(tom).IsA(admin)).To(BeFalse())
				Expect(Let(&tom).BeA(&admin)).To(Succeed())
				Expect(Whether(tom).IsA(admin)).To(BeTrue())
				Expect(Let(&tom).NotBeA(&admin)).To(Succeed())
				Expect(Whether(tom).IsA(admin)).To(BeFalse())
			})

			It("returns the temporary assignment status correctly", func() {
				Expect(Whether(tom).IsA(admin)).To(BeFalse())
				Expect(Let(&tom).BeATemp(admin)).To(Succeed())
				Expect(Whether(tom).IsA(admin)).To(BeTrue())
				Expect(Let(&tom).NotBeATemp(admin)).To(Succeed())
				Expect(Whether(tom).IsA(admin)).To(BeFalse())
			})
		})

		Describe("Permission of Role", func() {
			It("returns the assignment status correctly", func() {
				Expect(Whether(admin).Can(read_data)).To(BeFalse())
				Expect(Let(&admin).Can(&read_data)).To(Succeed())
				Expect(Whether(admin).Can(read_data)).To(BeTrue())
				Expect(Let(&admin).Cannot(&read_data)).To(Succeed())
				Expect(Whether(admin).Can(read_data)).To(BeFalse())
			})

			Context("Resource", func() {
				BeforeEach(func() {
					factory.Create(&resource)
					factory.Create(resource.NewPermissionForIt("read"))
				})

				It("can new the permission and returns the assignment status correctly", func() {
					Expect(Whether(resource).CanBe("read").By(admin)).To(BeFalse())
					Expect(Let(&resource).CanBe("read").By(&admin)).To(Succeed())
					Expect(Whether(resource).CanBe("read").By(admin)).To(BeTrue())
				})
			})
		})

		Describe("Permission of Subject", func() {
			Context("Role assignment first", func() {
				It("returns the nested assignment status correctly", func() {
					Expect(Let(&tom).BeA(&admin)).To(Succeed())
					Expect(Whether(admin).Can(read_data)).To(BeFalse())
					Expect(Whether(tom).Can(read_data)).To(BeFalse())

					Expect(Let(&admin).Can(&read_data)).To(Succeed())
					Expect(Whether(admin).Can(read_data)).To(BeTrue())
					Expect(Whether(&People{Name: "Tom"}).Can(read_data)).To(BeTrue())

					Expect(Let(&admin).Cannot(&read_data)).To(Succeed())
					Expect(Whether(admin).Can(read_data)).To(BeFalse())
					Expect(Whether(&People{Name: "Tom"}).Can(read_data)).To(BeFalse())
				})
			})

			Context("Role assignment last", func() {
				It("returns the nested assignment status correctly", func() {
					Expect(Whether(admin).Can(read_data)).To(BeFalse())
					Expect(Whether(tom).Can(read_data)).To(BeFalse())

					Expect(Let(&admin).Can(&read_data)).To(Succeed())
					Expect(Whether(admin).Can(read_data)).To(BeTrue())
					Expect(Whether(&People{Name: "Tom"}).Can(read_data)).To(BeFalse())
					Expect(Let(&tom).BeA(&admin)).To(Succeed())
					Expect(Whether(&People{Name: "Tom"}).Can(read_data)).To(BeTrue())

					Expect(Let(&tom).NotBeA(&admin)).To(Succeed())
					Expect(Whether(admin).Can(read_data)).To(BeTrue())
					Expect(Whether(&People{Name: "Tom"}).Can(read_data)).To(BeFalse())

					Expect(Let(&tom).BeA(&admin)).To(Succeed())
					Expect(Let(&admin).Cannot(&read_data)).To(Succeed())
					Expect(Whether(admin).Can(read_data)).To(BeFalse())
					Expect(Whether(&People{Name: "Tom"}).Can(read_data)).To(BeFalse())
				})
			})

			Context("Resource", func() {
				BeforeEach(func() {
					factory.Create(&resource)
					factory.Create(resource.NewPermissionForIt("read"))
				})

				It("can new the permission and returns the assignment status correctly", func() {
					Expect(Whether(resource).CanBe("read").By(admin)).To(BeFalse())
					Expect(Whether(resource).CanBe("read").By(tom)).To(BeFalse())
					Expect(Let(&resource).CanBe("read").By(&admin)).To(Succeed())
					Expect(Whether(resource).CanBe("read").By(admin)).To(BeTrue())
					Expect(Whether(resource).CanBe("read").By(tom)).To(BeFalse())
					Expect(Let(&tom).BeA(&admin)).To(Succeed())
					Expect(Whether(resource).CanBe("read").By(tom)).To(BeTrue())
				})
			})

			It("returns the temporary nested assignment status correctly", func() {
				Expect(Whether(admin).Can(read_data)).To(BeFalse())
				Expect(Whether(tom).Can(read_data)).To(BeFalse())

				Expect(Let(&admin).Can(&read_data)).To(Succeed())
				Expect(Whether(admin).Can(read_data)).To(BeTrue())
				Expect(Whether(tom).Can(read_data)).To(BeFalse())
				Expect(Let(&tom).BeATemp(admin)).To(Succeed())
				Expect(Whether(tom).Can(read_data)).To(BeTrue())

				Expect(Let(&tom).NotBeATemp(admin)).To(Succeed())
				Expect(Whether(admin).Can(read_data)).To(BeTrue())
				Expect(Whether(tom).Can(read_data)).To(BeFalse())

				Expect(Let(&tom).BeATemp(admin)).To(Succeed())
				Expect(Let(&admin).Cannot(&read_data)).To(Succeed())
				Expect(Whether(admin).Can(read_data)).To(BeFalse())

				// TODO: 临时角色赋予后，权限信息就无法变更了
				Expect(Whether(tom).Can(read_data)).To(BeTrue())
			})
		})
	})
}
