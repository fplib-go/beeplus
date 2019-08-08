package beeplus

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	Base
}

func (this *ErrorController) Error401() {
	this.Out(`
Error 401

    `).Die()
}
func (this *ErrorController) Error403() {
	this.Out(`
Error 403

    `).Die()
}
func (this *ErrorController) Error404() {
	this.Out(`
Error 404

    `).Die()
}
func (this *ErrorController) Error500() {
	this.Out(`
Error 500

    `).Die()
}
func (this *ErrorController) Error501() {
	this.Out(`
Error 501

    `).Die()

}
func (this *ErrorController) Error503() {
	this.Out(`
Error 503

    `).Die()

}
func (this *ErrorController) ErrorDb() {
	this.Out(`
Error DB

    `).Die()
}

func init() {
	// 注册错误处理函数
	beego.ErrorController(&ErrorController{})
}
