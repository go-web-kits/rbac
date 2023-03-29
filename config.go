package rbac

import (
	"github.com/go-web-kits/utils"
	"github.com/jinzhu/inflection"
)

type Config struct {
	Engine     Engine
	Role       string
	Permission string
	Subject    string // optional

	AutoDefineRole       bool
	AutoDefinePermission bool
}

var Configs = map[string]Config{}

func Init(configs map[string]Config) {
	for subject, c := range configs {
		c.Subject = subject
		Configs[subject] = c
		Configs[c.Role] = c
		Configs[c.Permission] = c
	}
}

func ConfigOf(obj interface{}) Config {
	key, ok := obj.(string)
	if !ok {
		key = utils.TypeNameOf(obj)
	}
	config, ok := Configs[key]
	if !ok {
		panic("Cannot Get Rbac Config Of " + key + ", Please Check Rbac Configuration!")
	}
	return config
}

func (c Config) Roles() string {
	return inflection.Plural(c.Role)
}

func (c Config) Permissions() string {
	return inflection.Plural(c.Permission)
}
