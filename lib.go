package beeplus

import (
	"github.com/fplib-go/fplib"
	"path"
	"strings"
)

type LibLoaderClass struct {
	Store map[string][]string
	Group map[string][]string
}

func (this *LibLoaderClass) Add(name string, libs ...string) {

	if this.Store == nil {
		this.Store = make(map[string][]string)
	}
	this.Store[name] = libs
}

func (this *LibLoaderClass) AddGroup(name string, libs ...string) {
	if this.Group == nil {
		this.Group = make(map[string][]string)
	}
	this.Group[name] = libs
}
func (this *LibLoaderClass) LoadCSS(name string, base ...string) string {
	_, out_css := this.GetLoadHTML(name, base...)
	return out_css
}
func (this *LibLoaderClass) LoadJS(name string, base ...string) string {
	out_js, _ := this.GetLoadHTML(name, base...)
	return out_js
}

func (this *LibLoaderClass) GetLoadHTML(name string, base ...string) (string, string) {
	p := ""
	if len(base) > 0 {
		p = base[0]
	}

	libs := []string{}
	names := fplib.Str.ToArr(name)
	for _, v := range names {
		G := this.Group[v]
		if len(G) > 0 {
			libs = append(libs, G...)
		} else {
			libs = append(libs, v)
		}

	}

	files := []string{}
	for _, v := range libs {
		L := this.Store[v]
		if len(L) > 0 {
			files = append(files, L...)
		}
	}
	files = fplib.Str.RemoveRepeatedFromArr(files)
	out_js, out_css := "", ""
	for _, v := range files {
		if strings.HasSuffix(v, ".css") {
			out_css += `<link href="` + path.Join(p, v) + `" rel="stylesheet">`
		} else if strings.HasSuffix(v, ".js") {
			out_js += `<script src="` + path.Join(p, v) + `"></script>`
		}

	}

	return out_js, out_css
}
