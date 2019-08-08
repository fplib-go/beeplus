package beeplus

import (
	"github.com/fplib-go/fplib"

	"github.com/astaxie/beego/orm"
	"github.com/tidwall/gjson"
)

type BaseUtils struct{}

// json转struct
func (this *BaseUtils) J2S(json interface{}, o interface{}, ps ...string) error {

	j := fplib.JSON.Obj(json).CamelString().Filter(ps...)
	fplib.Debug("j", j)
	err := fplib.JSON.ParseObj(j.Json, o)
	fplib.Debug("err", err)
	fplib.Debug("o", o)
	return err

	// rv := reflect.ValueOf(o).Elem()
	//
	// j.ForEach(func(key, value gjson.Result) bool {
	// 	k := key.String()
	// 	v := value.Value()
	// 	ku := fplib.Str.CamelString(k)
	//
	// 	fplib.Debug("v", v)
	// 	fplib.ShowType(v)
	// 	if f := rv.FieldByName(ku); f.CanSet() {
	// 		f.Set(reflect.ValueOf(v))
	// 	}
	//
	// 	return true // keep iterating
	// })

}

// struct转json
func (this *BaseUtils) S2J(o interface{}, ps ...string) *fplib.SJSON {
	return fplib.JSON.Obj(o).Filter(ps...).SnakeString()

}

// json转orm.Params
func (this *BaseUtils) J2ORM(o interface{}, ps ...string) orm.Params {
	j := fplib.JSON.Obj(o).Filter(ps...).GJSON()
	p := orm.Params{}
	j.ForEach(func(key, value gjson.Result) bool {
		k := key.String()
		v := value.Value()
		p[k] = v
		return true // keep iterating
	})
	return p
}
