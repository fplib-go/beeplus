package beeplus

import (
	"github.com/fplib-go/fplib"
	"fmt"

	"html/template"

	"github.com/astaxie/beego"
)

func init() {
	beego.AddFuncMap("LoadCSS", func(name string, c interface{}) template.HTML {
		if name == "" {
			name = "base"
		}

		var config RouterConfig
		for key, value := range c.(map[interface{}]interface{}) {
			strKey := fmt.Sprintf("%v", key)
			if strKey == "RouterConfig" {
				config = value.(RouterConfig)
				break
			}
		}
		html := config.LibLoader.LoadCSS(name, "/"+config.WebRoot)
		return template.HTML(html)
	})
	beego.AddFuncMap("LoadJS", func(name string, c interface{}) template.HTML {
		if name == "" {
			name = "base"
		}

		var config RouterConfig
		for key, value := range c.(map[interface{}]interface{}) {
			strKey := fmt.Sprintf("%v", key)
			if strKey == "RouterConfig" {
				config = value.(RouterConfig)
				break
			}
		}
		html := config.LibLoader.LoadJS(name, "/"+config.WebRoot)
		return template.HTML(html)
	})
	beego.AddFuncMap("json", func(o interface{}) template.JS {
		if s, ok := o.(string); ok {
			s = fplib.Trim(s)
			if fplib.Str.StartWith(s, "{") || fplib.Str.StartWith(s, "[") {
				if fplib.Str.EndWith(s, "}") || fplib.Str.EndWith(s, "]") || fplib.Str.EndWith(s, ";") || fplib.Str.EndWith(s, ",") {
					fplib.Debug("isjsonstr:", template.JS(s))
					return template.JS(s)
				}
			}
		}

		str := fplib.JSON.Stringify(o)
		if str == "null" {
			str = "{}"
		}
		return template.JS(str)
	})

	beego.AddFuncMap("site", func(key string) template.HTML {
		v := beego.AppConfig.String("site::" + key)
		return template.HTML(v)
	})

}
