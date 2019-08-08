package beeplus

import (
	"github.com/fplib-go/fplib"

	"github.com/astaxie/beego/cache"
)

var (
	Cache cache.Cache
)

func init() {
	var err error
	Cache, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		fplib.Error("Cache Error")
	}
}
