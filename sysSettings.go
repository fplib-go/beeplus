package beeplus

import (
	"github.com/fplib-go/fplib"
	"time"
)

var (
	Settings SettingsClass
)

type SettingsClass struct{}

func (this *SettingsClass) Set(key string, value string, remark ...string) bool {
	result := false
	db := DB()
	setting := SysSettings{Key: key}
	if _, _, err := db.ReadOrCreate(&setting, "Key"); err == nil {
		setting.Value = value
		if len(remark) > 0 {
			setting.Remark = remark[0]
		}
		if num, err := db.Update(&setting); err == nil {
			if num > 0 {
				result = true
			}
		}
	}

	return result
}

func (this *SettingsClass) Get(key string, defaults ...string) string {
	result := ""
	if len(defaults) > 0 {
		result = defaults[0]
	}

	db := DB()
	setting := SysSettings{Key: key}
	if created, _, err := db.ReadOrCreate(&setting, "Key"); err == nil {
		if created {
			// fmt.Println("New Insert an object. Id:", id)
		} else {
			// fmt.Println("Get an object. Id:", id)
			result = setting.Value
		}
	}
	return result
}

func (this *SettingsClass) GetInt(key string, defaults ...int) int {
	v := this.Get(key)
	result := 0
	if len(defaults) > 0 {
		result = defaults[0]
	}
	if v != "" {
		result = fplib.Int(v)
	}
	return result
}
func (this *SettingsClass) GetBool(key string, defaults ...bool) bool {
	v := this.Get(key)
	result := false
	if len(defaults) > 0 {
		result = defaults[0]
	}
	if v != "" {
		result = fplib.Bool(v)
	}
	return result

}

// 系统设置表
type SysSettings struct {
	Id       int       `orm:"auto;pk;description(ID)"`
	Key      string    `orm:"size(512);index;unique;description(key)"`
	Value    string    `orm:"null;type(text);description(value)"`
	CreateAt time.Time `orm:"auto_now_add;type(datetime);description(创建时间)"`
	UpdateAt time.Time `orm:"null;auto_now;type(datetime);description(更新时间)"`
	Remark   string    `orm:"null;type(text);description(备注)"`
}

func init() {
	RegTable(&SysSettings{})
	Settings = SettingsClass{}
}
