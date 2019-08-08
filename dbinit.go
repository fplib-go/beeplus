package beeplus

import (
	"github.com/fplib-go/fplib"
	// "time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db_DR         map[string]bool = map[string]bool{}
	db_installers []DBInstaller   = make([]DBInstaller, 0)
)

type DBInstaller interface {
	Install()
}

// 注册数据表模型
func RegTable(models ...interface{}) {
	orm.RegisterModel(models...)
}

// 加载数据库配置，注册数据库，传入配置文件字段标题
func MySQL(conf_names ...string) {
	if !db_DR["mysql"] {
		db_DR["mysql"] = true
		orm.RegisterDriver("mysql", orm.DRMySQL)
	}
	name := "MySQL"
	if len(conf_names) > 0 {
		name = conf_names[0]
	}

	dbname := beego.AppConfig.String(name + "::name")
	if dbname == "" {
		dbname = "default"
	}
	str := beego.AppConfig.String(name+"::user") + ":" +
		beego.AppConfig.String(name+"::password") + "@tcp(" +
		beego.AppConfig.String(name+"::host") + ":" +
		beego.AppConfig.String(name+"::port") + ")/" +
		beego.AppConfig.String(name+"::db") + "?charset=" +
		beego.AppConfig.String(name+"::charset")
	orm.RegisterDataBase(dbname, "mysql", str)

	auto, err := beego.AppConfig.Bool(name + "::autoinstall")
	if err == nil && auto {
		orm.RunSyncdb(dbname, false, false)
		DBInit(dbname)
	}

}

func DBInit(dbname string) {

	install_at := Settings.Get("db_init_at")
	if install_at == "" {
		for _, v := range db_installers {
			v.Install()
		}
		Settings.Set("db_init_at", fplib.Datetime())
	}

}

func DBInstall(installer DBInstaller) {
	db_installers = append(db_installers, installer)
}
