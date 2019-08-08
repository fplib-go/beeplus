package beeplus

import (
	"github.com/fplib-go/fplib"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/tidwall/gjson"
)

type ApiResult struct {
	Result bool        `json:"result";`
	Code   interface{} `json:"code";`
	Data   interface{} `json:"data";`
}

type Base struct {
	UID       int64
	Post_JSON gjson.Result
	Config    RouterConfig
	VData     *fplib.SJSON
	BaseUtils
	beego.Controller
}

// hook接口
type HookBeforeAction interface {
	Before()
}
type HookAfterAction interface {
	After()
}

// 设置配置
func (this *Base) SetConfig(conf RouterConfig) {
	this.Config = conf

}

func (this *Base) Prepare() {
	// page start time
	this.Data["PageStartTime"] = time.Now()
	this.Data["RouterConfig"] = this.Config
	this.Data["SelfUrl"] = this.URL()
	this.Data["SubDir"] = this.Config.WebRoot

	this.VData = fplib.JSON.Obj()

	this.UID = 0
	if _uid, ok := this.Data["_uid"]; ok {
		if __uid, ok := _uid.(int64); ok {
			this.UID = __uid
		}
	}

	if app, ok := this.AppController.(HookBeforeAction); ok {
		app.Before()
	}
}
func (this *Base) Finish() {
	// page start time
	this.Data["PageEndTime"] = time.Now()

	if app, ok := this.AppController.(HookAfterAction); ok {
		app.After()
	}
}

// api正常输出
func (this *Base) Ok(obj ...interface{}) {
	this.ApiOut(true, obj...)
}

// api错误输出
func (this *Base) Err(obj ...interface{}) {
	this.ApiOut(false, obj...)
}

// api输出
func (this *Base) ApiOut(result bool, obj ...interface{}) {
	var mystruct ApiResult
	switch len(obj) {
	case 0:
		mystruct = ApiResult{
			Result: result,
			Code:   "",
			Data:   "",
		}
	case 1:

		mystruct = ApiResult{
			Result: result,
			Code:   "",
			Data:   fplib.JSON.ToObj(obj[0]),
		}

	case 2:
		mystruct = ApiResult{
			Result: result,
			Code:   fplib.JSON.ToObj(obj[1]),
			Data:   fplib.JSON.ToObj(obj[0]),
		}

	default:
		mystruct = ApiResult{
			Result: false,
			Code:   "",
			Data:   "Server Error",
		}

	}

	this.Data["json"] = &mystruct
	this.ServeJSON()
	this.Die()
}

// 字符串输出
func (this *Base) Out(data string) *Base {
	this.Ctx.WriteString(data)
	return this
}

// 结束
func (this *Base) Die() {
	this.Finish()
	this.StopRun()
}

// 跳转网页
func (this *Base) Goto(url string) {
	this.Redirect(url, 302)
	this.Die()
	// this.Ctx.Redirect(302, url)
}

// 调用模板
func (this *Base) View(files ...string) {
	if !fplib.Empty(this.Config.ViewRoot) {
		this.TplPrefix = fplib.Trim(this.Config.ViewRoot, "/") + "/"
	}

	switch len(files) {
	case 0:
		c, a := this.GetControllerAndAction()
		this.TplName = strings.ToLower(c) + "/" + strings.ToLower(a) + "." + this.Config.GetString("TplExt", "html")
		this.Layout = this.getLayoutPath("")
	case 1:
		this.TplName = files[0]
		this.Layout = this.getLayoutPath("")
	case 2:
		this.TplName = files[0]
		if files[1] != "" {
			this.Layout = this.getLayoutPath(files[1])
		}
	}

	b, err := beego.AppConfig.Bool("autorender")

	if (err == nil) && (!b) {
		this.Render()
	}
}

// 获取访问地址
func (this *Base) URL(urls ...string) string {
	url := ""
	if len(urls) > 0 {
		url = path.Join(urls...)
	}
	if url == "" || url == "." {
		url = this.Ctx.Input.URI()
	} else {
		url = this.MakeURL(url)
	}

	return url
}

// 获取带前缀的url
func (this *Base) MakeURL(url string) string {
	return path.Join("/", this.Config.WebRoot, url)
}

// 获取layout路径
func (this *Base) getLayoutPath(file string) string {
	if file == "" {
		file = "index." + this.Config.GetString("TplExt", "html")
	}
	return fplib.Trim(path.Join("/", this.Config.ViewRoot, "layout", file), "/")
}

// 获取post上来的json数据
func (this *Base) GetJSON(keys ...string) gjson.Result {
	if !this.Post_JSON.Exists() {
		this.Post_JSON = fplib.JSON.Parse(string(this.Ctx.Input.RequestBody))
	}
	if len(keys) > 0 {
		key := strings.Join(keys, ".")
		return this.Post_JSON.Get(key)
	}
	return this.Post_JSON

}

// 获取post上来的表单数据
func (this *Base) GetMap(key string) gjson.Result {
	var keys map[string]string
	keys = make(map[string]string, 0)
	this.Ctx.Input.Bind(&keys, key)

	// return keys
	return fplib.JSON.Obj(keys).GJSON()
}

// 设置gjson的默认值
func (this *Base) Default(json gjson.Result, defaults ...interface{}) gjson.Result {
	make := false
	if json.Exists() {
		v := json.Value()
		if fplib.Empty(v) {
			make = true
		}
	} else {
		make = true
	}

	if make {
		if len(defaults) > 0 {
			def := defaults[0]
			// r := gjson.Get("{\"a\":\""+def+"\"}", "a")
			json = gjson.Parse(fplib.JSON.Stringify(def))
		}
	}
	return json

}

// 获取变量，可设置默认值
func (this *Base) V(key string, defaults ...interface{}) gjson.Result {
	// 如果没解析过变量，解析所有变量到VData
	if this.VData.IsEmpty() {
		this.ParseAllVData()
	}
	return this.Default(this.VData.Get(key), defaults...)
}

// 解析所有变量到VData
func (this *Base) ParseAllVData() {

	// 处理Params
	this.VData.Parse(this.Ctx.Input.Params())

	// 处理post数据
	arr := this.Input()

	for k, v := range arr {
		if strings.Contains(k, "[]") {
			index := strings.IndexAny(k, "[")
			kk := k[:index]
			this.VData.Set(kk, v)
		} else if strings.Contains(k, "[") && strings.Contains(k, "]") {
			index := strings.IndexAny(k, "[")
			kk := k[:index]
			vs := make(map[string]string, 0)
			this.Ctx.Input.Bind(&vs, kk)
			this.VData.Set(kk, vs)
		} else {

			if vv, ok := interface{}(v).([]string); ok && len(vv) == 1 {
				this.VData.Set(k, vv[0])
			} else {
				this.VData.Set(k, v)
			}

		}
		this.VData.Set("\\:dataType", "form")
	}

	// 是否是json格式的body
	if strings.Contains(this.Ctx.Input.Header("Content-Type"), "json") {
		j := this.GetJSON()
		j.ForEach(func(key, value gjson.Result) bool {
			this.VData.Set(key.String(), value.Value())
			return true // keep iterating
		})
		this.VData.Set("\\:dataType", "json")
	}

}

// 判断变量是否为空，如果为空直接输出api错误
func (this *Base) Has(v interface{}, msgs ...string) *Base {

	if fplib.Empty(v) {
		msg := "参数错误"
		if len(msgs) > 0 {
			msg = strings.Join(msgs, " ")
		}
		this.Err(msg)
	}
	return this
}

func (this *Base) HasAll(vs ...interface{}) *Base {
	for _, v := range vs {
		this.Has(v)
	}
	return this
}

// 带模板的页面错误提示
func (this *Base) PageErr(msg string) {
	this.Data["ErrorMsg"] = msg
	this.View("error.html")
	this.Die()

}

// func (this *Base) VueData(key string, val interface{}) {
// 	if _, ok := this.Data["_VueData"]; !ok {
// 		this.Data["_VueData"] = make(map[string]interface{}, 0)
// 	}
// 	this.Data["_VueData"][key] = val
// }
