package beeplus

import (
	"github.com/fplib-go/fplib"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

type RouterIndex struct {
	Key        string
	Param      string
	Action     string
	Controller interface{}
	Remark     string
}

type BaseControllerInterface interface {
	beego.ControllerInterface
	SetConfig(RouterConfig)
}

type GoIndexController struct {
	beego.Controller
}

func (this *GoIndexController) GoIndex() {
	url := ""
	if rp, ok := this.Data["RouterPattern"].(string); ok {
		url = path.Join(rp, "index")
	} else {
		url = path.Join(this.Ctx.Input.URL(), "index")
	}
	this.Ctx.Redirect(301, url)
}

// 注册路由
func Router(base string, routers []RouterIndex, config ...RouterConfig) {
	for _, v := range routers {
		if vc, ok := v.Controller.(BaseControllerInterface); ok {
			// 为每一个控制器设置config
			if len(config) > 0 {
				vc.SetConfig(config[0])
			} else {
				vc.SetConfig(RouterConfig{})
			}

			// 处理变量
			key := fplib.Trim(v.Key)
			param := fplib.Default(v.Param, "*").(string)
			action := fplib.Default(v.Action, "index").(string)
			actions := fplib.Str.ToArr(fplib.Str.SnakeString(action))

			// fplib.Debug("action", action)
			// fplib.Debug("v.action", v.action)

			// 增加首页默认路由index
			if fplib.Bool(Options["autoIndex"]) {
				if !fplib.Str.In_Array(actions, "index") {
					actions = append(actions, "index")
				}
				// 增加没有action情况下的默认路由
				indexurl := path.Join(base, key)
				beego.Router(indexurl, vc, "*:Index")
				// beego.Router(indexurl, &GoIndexController{}, "*:GoIndex")

			}

			// 循环所有action
			for _, a := range actions {
				url := path.Join(base, key, a, param)
				api := "*:" + strings.Title(a)
				// fplib.Debug("url", url, "api", api)
				beego.Router(url, vc, api)
			}
		}

	}
}

// 注册静态lib路径
func Lib(base, dir string) {
	beego.SetStaticPath(base, dir)
}

// 以web根目录启动一个web框架
func WWWROOT(config RouterConfig) {

	webBase := config.WebRoot
	srcBase := config.SrcRoot
	webroot := path.Join("/", webBase)

	page := config.PageIndex
	api := config.ApiIndex

	Router(webroot, page, config)
	Router(path.Join(webroot, "api"), api, config)

	Lib(path.Join(webroot, "data"), path.Join(srcBase, "data"))
	Lib(path.Join(webroot, "images"), path.Join(srcBase, "images"))
	Lib(path.Join(webroot, "lib"), path.Join(srcBase, "lib"))

}
