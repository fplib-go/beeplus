package beeplus

import (
	"github.com/astaxie/beego"
)

// 注册Hook
func Hook(hash, position string, hook beego.FilterFunc, params ...bool) {
	hook_position := beego.BeforeStatic
	switch position {
	case "before":
		hook_position = beego.BeforeExec
	case "after":
		hook_position = beego.AfterExec
	case "BeforeStatic":
		hook_position = beego.BeforeStatic
	case "BeforeRouter":
		hook_position = beego.BeforeRouter
	case "BeforeExec":
		hook_position = beego.BeforeExec
	case "AfterExec":
		hook_position = beego.AfterExec
	case "FinishRouter":
		hook_position = beego.FinishRouter
	}
	beego.InsertFilter(hash, hook_position, hook, params...)
}
