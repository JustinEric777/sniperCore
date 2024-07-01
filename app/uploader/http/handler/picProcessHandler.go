package handler

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sniperCore/app/uploader/http/request"
	"github.com/sniperCore/app/uploader/service"
	"github.com/sniperCore/core/config"
	"github.com/sniperCore/untils"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	picProcessService *service.PictureProcessService
)

func PictureProcess(ctx *gin.Context) {
	param := &request.PicProcessParam{}
	err := param.PicProcessValidate(ctx)
	if err != nil {
		untils.JsonError(ctx, 400, "params check error，err = "+err.Error())
		ctx.Abort()
		return
	}

	file, err := param.Content.Open()
	if err != nil {
		untils.JsonError(ctx, 400, "get file error，err = "+err.Error())
		ctx.Abort()
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		untils.JsonError(ctx, 500, "read file's data error，err = "+err.Error())
		ctx.Abort()
		return
	}
	key, _ := GetPicturePath(param.Scene, param.CloudFileName)
	input, err := generatePicture(string(content), key)
	if err != nil {
		untils.JsonError(ctx, 500, "file generate error，err = "+err.Error())
		ctx.Abort()
		return
	}
	if input == "" {
		input = param.FileName
	}

	var files []map[string]string
	fileInfo := make(map[string]string)
	fileInfo["key"] = key
	fileInfo["filename"] = input
	files = append(files, fileInfo)
	dir := config.GetAppPath() + "/storage/picture/"
	for _, action := range param.Commands {
		output, key := action.Output, action.Output
		if IsCloud(param.IsCloud) {
			key, err = GetPicturePath(param.Scene, action.Output)
			if err != nil {
				untils.JsonError(ctx, 500, "file's key generate error，err = "+err.Error())
				ctx.Abort()
				return
			}
			output = dir + key
		}

		dir := filepath.Dir(output)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				untils.JsonError(ctx, 500, "output dir create error，err = "+err.Error())
				ctx.Abort()
				return
			}
		}
		if err := picProcessService.Process(input, output, action.Cmd); err != nil {
			untils.JsonError(ctx, 500, "picture process error，err = "+err.Error())
			ctx.Abort()
			return
		}

		fileInfo = make(map[string]string)
		fileInfo["key"] = key
		fileInfo["filename"] = output
		files = append(files, fileInfo)
	}

	untils.JsonSuccess(ctx, files)
}

func generatePicture(content, filename string) (string, error) {
	if content != "" {
		baseDir := config.GetAppPath() + "/storage/picture"
		filename = baseDir + "/" + filename
		dir := filepath.Dir(filename)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return "", err
			}
		}
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		if err != nil {
			return "", err
		}
		defer file.Close()

		write := bufio.NewWriter(file)
		_, err = write.WriteString(content)
		if err != nil {
			return "", err
		}

		return filename, nil
	} else {
		return "", nil
	}
}

func IsCloud(isCloud bool) bool {
	if !isCloud {
		return false
	}

	return true
}

func GetPicturePath(scene, fileName string) (string, error) {
	configure := config.Conf
	paths := make(map[string]string)
	if err := configure.UnmarshalKey("picPaths", &paths); err != nil {
		return "", err
	}

	subPath := time.Now().Format("2006-01-02")
	if path, existed := paths[scene]; existed {
		return fmt.Sprintf("%s", path) + "/" + subPath + "/" + fileName, nil
	}

	return "", errors.New("not found scene")
}
