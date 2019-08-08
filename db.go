package beeplus

import (
	"github.com/fplib-go/fplib"
	"errors"
	"reflect"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// 数据表的基础
type TableBase struct {
	Id       int64     `orm:"auto;pk;description(ID)"`
	Status   int64     `orm:"default(1);index;description(状态)"`
	CreateAt time.Time `orm:"auto_now_add;type(datetime);description(创建时间)"`
	UpdateAt time.Time `orm:"null;auto_now;type(datetime);description(更新时间)"`
	DeleteAt time.Time `orm:"null;type(datetime);description(删除时间)"`
	Remark   string    `orm:"null;type(text);description(备注)"`
}

// 获取数据库对象
func DB(db ...string) *DBClass {
	dbc := DBClass{}
	return dbc.Init(db...)
}

type DBClass struct {
	Orm   orm.Ormer
	Seter orm.RawSeter
	Count int64
}

// 初始化
func (this *DBClass) Init(db ...string) *DBClass {
	this.Orm = orm.NewOrm()
	if len(db) > 0 {
		this.Orm.Using(db[0])
	}
	return this
}

// 开启debug
func (this *DBClass) Debug(b ...bool) {
	e := true
	if len(b) > 0 {
		e = b[0]
	}
	orm.Debug = e
}

// 读取
func (this *DBClass) Read(t interface{}, p ...string) error {
	return this.Orm.Read(t, p...)
}

// 读取，没有就创建
func (this *DBClass) ReadOrCreate(t interface{}, p string, ps ...string) (bool, int64, error) {
	return this.Orm.ReadOrCreate(t, p, ps...)
}

// 插入
func (this *DBClass) Insert(t interface{}) (int64, error) {
	return this.Orm.Insert(t)
}

// 插入新，并使记录保持可用 status=1 delete_at=null
func (this *DBClass) I(t interface{}) (int64, error) {
	this.SetEnable(t)
	return this.Orm.Insert(t)
}

// 批量插入
func (this *DBClass) InsertMulti(c int, t interface{}) (int64, error) {
	return this.Orm.InsertMulti(c, t)
}

// 更新
func (this *DBClass) Update(t interface{}, p ...string) (int64, error) {
	return this.Orm.Update(t, p...)
}

// 更新，并使记录保持可用 status=1 delete_at=null
func (this *DBClass) U(t interface{}, p ...string) (int64, error) {
	this.SetEnable(t)
	return this.Orm.Update(t, p...)
}

// 删除
func (this *DBClass) Delete(t interface{}, p ...string) (int64, error) {
	return this.Orm.Delete(t, p...)
}

// 禁用，软删除功能，直接写入数据库
func (this *DBClass) Disable(t interface{}, p ...string) (int64, error) {
	this.SetDisable(t)
	return this.Orm.Update(t, p...)
}

// 设置禁用状态，不写入数据库
func (this *DBClass) SetDisable(t interface{}) {
	v := reflect.ValueOf(t).Elem()
	if f := v.FieldByName("Status"); f.CanSet() {
		f.SetInt(0)
	}
	if f := v.FieldByName("DeleteAt"); f.CanSet() {
		f.Set(reflect.ValueOf(time.Now()))
	}
}

// 设置启用状态，直接写入数据库
func (this *DBClass) Enable(t interface{}, p ...string) (int64, error) {
	this.SetEnable(t)
	return this.Orm.Update(t, p...)
}

// 设置启用状态，不写入数据库
func (this *DBClass) SetEnable(t interface{}) {
	v := reflect.ValueOf(t).Elem()
	if f := v.FieldByName("Status"); f.CanSet() {
		f.SetInt(1)
	}
	if f := v.FieldByName("DeleteAt"); f.CanSet() {
		f.Set(reflect.ValueOf(time.Time{}))
	}
}

// 执行sql
func (this *DBClass) Raw(sql string, p ...interface{}) orm.RawSeter {
	this.Seter = this.Orm.Raw(sql, p...)
	return this.Seter
}

// 通过表名返回QueryTable
func (this *DBClass) Table(t interface{}) orm.QuerySeter {
	return this.Orm.QueryTable(t)
}

// 通过查询字段返回QueryBuilder
func (this *DBClass) Select(p ...string) orm.QueryBuilder {
	dbtype := "mysql"
	if this.Orm != nil {
		dr := this.Orm.Driver()
		switch dr.Type() {
		case orm.DRMySQL:
			dbtype = "mysql"
		case orm.DRSqlite:
			dbtype = "sqlite3"
		case orm.DRPostgres:
			dbtype = "postgres"
		}
	}
	qb, _ := orm.NewQueryBuilder(dbtype)
	return qb.Select(p...)
}

// 执行查询语句，返回结果数组
func (this *DBClass) Q(sql string, p ...interface{}) []orm.Params {
	this.Raw(sql, p...)
	return this.Arr()
}

// 从Raw查询中解析结果到数组
func (this *DBClass) Arr(seters ...interface{}) []orm.Params {
	this.Count = 0
	var maps []orm.Params
	var num int64
	var err error
	if len(seters) > 0 {
		seter := seters[0]
		switch v := seter.(type) {
		case orm.QuerySeter:
			num, err = v.Values(&maps)
		case orm.RawSeter:
			num, err = v.Values(&maps)
		case string:
			if len(seters) == 1 {
				num, err = this.Orm.Raw(v).Values(&maps)
			} else {
				ps := seters[1:]
				num, err = this.Orm.Raw(v, ps...).Values(&maps)
			}
		default:
			num = 0
			err = errors.New("Seter Error")
		}
	} else {
		num, err = this.Seter.Values(&maps)
	}

	if err != nil {
		fplib.Error(err)
	}
	this.Count = num
	return maps
}

// 从Raw查询中解析一条结果到数组
func (this *DBClass) One(seters ...interface{}) orm.Params {
	maps := this.Arr(seters...)
	m := make(orm.Params)
	if len(maps) > 0 {
		m = maps[0]
		this.Count = 1
	}
	return m
	// this.Count = 0
	// var maps []orm.Params
	// m := make(orm.Params)
	// num, err := this.Seter.Values(&maps)
	// if err == nil && num > 0 {
	// 	m = maps[0]
	// 	this.Count = 1
	// }
	// return m
}

// 设置Raw的参数
func (this *DBClass) P(p ...interface{}) orm.RawSeter {
	this.Seter = this.Seter.SetArgs(p...)
	return this.Seter
}

func (this *DBClass) Begin() {
	this.Orm.Begin()
}
func (this *DBClass) Rollback() {
	this.Orm.Rollback()
}
func (this *DBClass) Commit() {
	this.Orm.Commit()
}

// 在闭包函数中启用事务处理，随时可以panic,将自动回滚
func (this *DBClass) T(f func()) (Rerr error) {
	Rerr = this.Orm.Begin()
	if Rerr == nil {
		defer func() {
			if p := recover(); p != nil {
				this.Orm.Rollback()
				switch v := p.(type) {
				case string:
					Rerr = errors.New(v)
				case error:
					Rerr = v
				default:
					Rerr = errors.New("DB Error")
				}

			} else {
				Rerr = this.Orm.Commit()
			}
		}()
		f()
	}
	return

}
