package domain

import "mime/multipart"

type UploadFile struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	AuthorId  int64  `json:"AuthorId"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"createdAt"`
}

type UploadForm struct {
	File *multipart.FileHeader `form:"file" json:"file"`
	Json UploadFile            `form:"json" json:"json"`
}
