package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type UploadParam struct {
	CloudFileName string                `form:"cloud_filename" validate:"required"`
	Content       *multipart.FileHeader `form:"content" validate:"required"`
	AppCode       string                `form:"app_code" validate:"required"`
	Scene         string                `form:"scene" validate:"required"`
	Cloud         string                `form:"cloud"`
}

func (param *UploadParam) UploadValidate(ctx *gin.Context) error {
	if err := ctx.Bind(param); err != nil {
		return err
	}
	cloud := getAppCloud(param.AppCode)
	if cloud == "" {
		return errors.New("not supported app")
	}
	param.Cloud = cloud

	return nil
}

type UploadLocalParam struct {
	CloudFileName string `form:"cloud_filename" validate:"required"`
	FileName      string `form:"filename" validate:"required"`
	AppCode       string `form:"app_code" validate:"required"`
	Scene         string `form:"scene" validate:"required"`
	Cloud         string `form:"cloud"`
}

func (param *UploadLocalParam) UploadLocalValidate(ctx *gin.Context) error {
	if err := ctx.ShouldBind(param); err != nil {
		return err
	}
	cloud := getAppCloud(param.AppCode)
	if cloud == "" {
		return errors.New("not supported app")
	}
	param.Cloud = cloud

	return nil
}

type RemoveObjParam struct {
	Key     string `form:"key" validate:"required"`
	AppCode string `form:"app_code" validate:"required"`
	Cloud   string `form:"cloud"`
}

func (param *RemoveObjParam) RemoveObjValidate(ctx *gin.Context) error {
	if err := ctx.ShouldBind(param); err != nil {
		return err
	}
	cloud := getAppCloud(param.AppCode)
	if cloud == "" {
		return errors.New("not supported app")
	}
	param.Cloud = cloud

	return nil
}

type GetObjectStreamParam struct {
	Key     string `form:"key" validate:"required"`
	AppCode string `form:"app_code" validate:"required"`
	Cloud   string `form:"cloud"`
}

func (param *GetObjectStreamParam) GetObjectStreamValidate(ctx *gin.Context) error {
	if err := ctx.ShouldBind(param); err != nil {
		return err
	}
	cloud := getAppCloud(param.AppCode)
	if cloud == "" {
		return errors.New("not supported app")
	}
	param.Cloud = cloud

	return nil
}

type GetObjectUrlParam struct {
	Key        string `form:"key" validate:"required"`
	AppCode    string `form:"app_code" validate:"required"`
	ExpireTime int    `form:"expire_time" validate:"required"`
	Cloud      string `form:"cloud"`
}

func (param *GetObjectUrlParam) GetObjectUrlValidate(ctx *gin.Context) error {
	if err := ctx.ShouldBind(param); err != nil {
		return err
	}
	cloud := getAppCloud(param.AppCode)
	if cloud == "" {
		return errors.New("not supported app")
	}
	param.Cloud = cloud

	return nil
}

type GetObjectToLocalParam struct {
	Key      string `form:"key" validate:"required"`
	AppCode  string `form:"app_code" validate:"required"`
	Filename string `form:"filename" validate:"required"`
	Cloud    string `form:"cloud"`
}

func (param *GetObjectToLocalParam) GetObjectToLocalValidate(ctx *gin.Context) error {
	if err := ctx.ShouldBind(param); err != nil {
		return err
	}
	cloud := getAppCloud(param.AppCode)
	if cloud == "" {
		return errors.New("not supported app")
	}
	param.Cloud = cloud

	return nil
}

func getAppCloud(appCode string) string {
	var cloud string

	switch appCode {
	case "dq":
		cloud = "dq_pic_cloud"
	case "dy":
		cloud = "dy_pic_cloud"
	}

	return cloud
}
