package rbac_dbx_test

import (
	"fmt"
	"testing"

	"github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/rbac"
	"github.com/go-web-kits/rbac/rbac_dbx"
	. "github.com/go-web-kits/rbac/rbac_dbx/test"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/utils/project"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRbacDBx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rbac(DBx) Suite")
}

var _ = BeforeSuite(func() {
	dbx.Client, _ = getPostgresDB()
	Boot().Migrate(&People{}, &PeopleRole{}, &PeoplePermission{}, &PeopleResource{})

	Init(map[string]Config{
		"People": {Engine: rbac_dbx.Engine{}, Role: "PeopleRole", Permission: "PeoplePermission"}})
})

var _ = AfterSuite(func() {
	ShutApp()
})

// ======

// TODO: extra
type DatabaseConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"user_name"`
	Database     string `json:"database"`
	TestDatabase string `json:"test_database"`
	Password     string `json:"password"`
	DisableSSL   bool   `json:"disable_ssl"`
}

func (c *DatabaseConfig) ConnectString() string {
	db := c.Database
	if project.OnTest() {
		db = c.TestDatabase
	}
	fmt.Println("* Connecting Database\n    using `" + db + "`")

	s := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		c.Host, c.Port, c.Username, db, c.Password,
	)
	if c.DisableSSL {
		s += " sslmode=disable"
	}

	return s
}

func getPostgresDB() (db *gorm.DB, err error) {
	c := DatabaseConfig{
		Host: "localhost", Port: 5432, Database: "firmware_services_test",
		Username: "postgres", Password: "123456", DisableSSL: true,
	}
	db, err = gorm.Open("postgres", c.ConnectString())
	if err != nil {
		fmt.Println("    ERROR: " + err.Error())
		fmt.Println("    Host: " + c.Host)
		return nil, err
	}
	fmt.Println("    Connected!")
	return db, err
}
