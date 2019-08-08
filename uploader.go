package beeplus

import (
	"github.com/fplib-go/fplib"
	"mime/multipart"
)

type UploadResult struct {
	Result bool        `json:"result";`
	Code   interface{} `json:"code";`
	Data   interface{} `json:"data";`
	// 兼容ckeditor的输出
	Uploaded bool   `json:"uploaded";`
	Url      string `json:"url";`
}

type fileInfo struct {
	has       bool
	fieldName string
	file      multipart.File
	header    *multipart.FileHeader
}

type Uploader struct {
	Base
}

func (this *Uploader) GetUploadFile(baseDir string) {
	names := []string{
		"upload",
		"file",
		"files",
		"image",
		"images",
	}
	info := fileInfo{}
	for _, v := range names {

		var err error
		info.file, info.header, err = this.GetFile(v)
		if err == nil {
			defer info.file.Close()
			info.has = true
			info.fieldName = v

			break
		}
	}

	if info.has {
		// targetDir := path.Join(baseDir, info.header.Filename)
		this.SaveToFile(
			info.fieldName,
			"static/upload/"+info.header.Filename,
		) // 保存位置在 static/upload, 没有文件夹要先创建

	} else {
		this.Err("遇到错误")
	}

}

// api正常输出
func (this *Uploader) Ok(obj interface{}) {
	o := fplib.JSON.Obj(obj)
	mystruct := UploadResult{
		Result:   true,
		Uploaded: true,
		Data:     o,
		Url:      o.Get("url").String(),
	}
	this.Data["json"] = &mystruct
	this.ServeJSON()
	this.Die()
}

// api错误输出
func (this *Uploader) Err(obj interface{}) {
	mystruct := UploadResult{
		Result: false,
		Data:   fplib.JSON.ToObj(obj),
	}
	this.Data["json"] = &mystruct
	this.ServeJSON()
	this.Die()
}
