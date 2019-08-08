package beeplus

import (
	"github.com/fplib-go/fplib"
)

type RouterConfig struct {
	WebRoot   string
	SrcRoot   string
	ViewRoot  string
	LibLoader *LibLoaderClass
	PageIndex []RouterIndex
	ApiIndex  []RouterIndex
	fplib.Conf
}

func NewRouterConfig(p ...string) RouterConfig {
	webroot := "/"
	srcroot := "/"
	viewroot := "views"
	switch len(p) {
	case 1:
		webroot = p[0]
		srcroot = p[0]
		viewroot = p[0]
	case 2:
		webroot = p[0]
		srcroot = p[1]
		viewroot = p[1]
	case 3:
		webroot = p[0]
		srcroot = p[1]
		viewroot = p[2]
	}
	return RouterConfig{
		WebRoot:   webroot,
		SrcRoot:   srcroot,
		ViewRoot:  viewroot,
		LibLoader: &LibLoaderClass{},
		PageIndex: make([]RouterIndex, 0),
		ApiIndex:  make([]RouterIndex, 0),
	}
}
