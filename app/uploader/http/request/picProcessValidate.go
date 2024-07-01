package request

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type PicProcessParam struct {
	FileName      string                `form:"filename"`
	CloudFileName string                `form:"cloud_filename"`
	Content       *multipart.FileHeader `form:"content"`
	AppCode       string                `form:"app_code" validate:"required" default:"true"`
	IsCloud       bool                  `form:"is_cloud" validate:"required"`
	Actions       string                `form:"actions" validate:"required"`
	Commands      []*Action             `form:"commands" `
	Scene         string                `form:"scene" validate:"required"`
	Cloud         string                `form:"cloud"`
}

type Action struct {
	Cmd    string `json:"cmd"`
	Output string `json:"output"`
}

func (param *PicProcessParam) PicProcessValidate(ctx *gin.Context) error {
	if err := ctx.Bind(param); err != nil {
		return err
	}
	cloud := getAppCloud(param.AppCode)
	if cloud == "" {
		return errors.New("not supported app")
	}
	param.Cloud = cloud
	if param.FileName == "" && param.Content == nil {
		return errors.New("filename and content is empty")
	}
	var actions []*Action
	err := json.Unmarshal([]byte(param.Actions), &actions)
	if err != nil || len(actions) < 1 {
		return errors.New("actions parse exception")
	}
	param.Commands = actions

	return nil
}
