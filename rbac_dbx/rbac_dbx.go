package rbac_dbx

import "github.com/go-web-kits/dbx"

type Engine struct{}

func getOpts(args ...interface{}) []dbx.Opt {
	opts := []dbx.Opt{}
	for _, arg := range args {
		opts = append(opts, arg.(dbx.Opt))
	}
	return opts
}
