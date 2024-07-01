package service

import "C"
import (
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"os/exec"
	"strconv"
	"strings"
)

type PictureProcessService struct {
}

// 单条处理原因因为多个工具柔和
func (ps *PictureProcessService) Process(input, output, cmdStr string) error {
	cmdArr := strings.Split(cmdStr, "&")
	for _, cmd := range cmdArr {
		op := strings.Split(cmd, ":")
		switch op[0] {
		case "exif":
			if err := ps.Exif(input, output); err != nil {
				return err
			}
		case "resize":
			params := parseParams(op[1])
			width, _ := strconv.Atoi(params["width"].(string))
			height, _ := strconv.Atoi(params["height"].(string))

			if err := ps.Resize(input, output, width, height); err != nil {
				return err
			}
		case "crop":
			params := parseParams(op[1])
			width, _ := strconv.Atoi(params["width"].(string))
			height, _ := strconv.Atoi(params["height"].(string))
			x, _ := strconv.Atoi(params["x"].(string))
			y, _ := strconv.Atoi(params["y"].(string))
			gravity := params["gravity"].(string)

			if err := ps.Crop(input, output, width, height, x, y, gravity); err != nil {
				return err
			}
		case "compress":
			params := parseParams(op[1])
			quality, _ := strconv.Atoi(params["quality"].(string))
			encoder, _ := params["encoder"].(string)

			if err := ps.Compress(input, output, encoder, quality); err != nil {
				return err
			}
		case "convert":
			params := strings.Split(op[1], "=")

			if _, err := ps.Convert(input, output, params...); err != nil {
				return err
			}
		}
		input = output
	}

	return nil
}

func parseParams(paramStr string) map[string]interface{} {
	paramsMap := make(map[string]interface{})

	attrs := strings.Split(paramStr, ",")
	for _, attr := range attrs {
		params := strings.Split(attr, "=")
		paramsMap[params[0]] = params[1]
	}

	return paramsMap
}

// magick -size 16000x16000 -depth 8 -resize 640x480 image.rgb image.png
func (ps *PictureProcessService) Resize(input, output string, width, height int) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImage(input); err != nil {
		return err
	}

	if err := mw.ResizeImage(uint(width), uint(height), imagick.FILTER_LANCZOS); err != nil {
		return err
	}

	if err := mw.WriteImage(output); err != nil {
		return err
	}

	return nil
}

// magick image.png -gravity Center -region 10x10-40+20 -negate output.png
func (ps *PictureProcessService) Crop(input, output string, width, height, x, y int, gravity string) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImage(input); err != nil {
		return err
	}

	var gravityType imagick.GravityType
	if gravity != "" {
		switch gravity {
		case "undefined":
			gravityType = imagick.GRAVITY_CENTER
		case "center":
			gravityType = imagick.GRAVITY_UNDEFINED
		case "forget":
			gravityType = imagick.GRAVITY_FORGET
		case "north_west":
			gravityType = imagick.GRAVITY_NORTH_WEST
		case "north_east":
			gravityType = imagick.GRAVITY_NORTH_EAST
		case "south_west":
			gravityType = imagick.GRAVITY_SOUTH_WEST
		case "south_east":
			gravityType = imagick.GRAVITY_SOUTH_EAST
		case "north":
			gravityType = imagick.GRAVITY_NORTH
		case "west":
			gravityType = imagick.GRAVITY_WEST
		case "east":
			gravityType = imagick.GRAVITY_EAST
		case "south":
			gravityType = imagick.GRAVITY_SOUTH
		}
		if err := mw.SetGravity(gravityType); err != nil {
			return err
		}
	}
	if err := mw.CropImage(uint(width), uint(height), x, y); err != nil {
		return err
	}

	if err := mw.WriteImage(output); err != nil {
		return err
	}

	return nil
}

func (ps *PictureProcessService) Compress(input, output, encoder string, quality int) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImage(input); err != nil {
		return err
	}

	switch encoder {
	case "imagemagick":
		if err := mw.SetImageCompressionQuality(uint(quality)); err != nil {
			return err
		}
		if err := mw.WriteImage(output); err != nil {
			return err
		}
	case "guetzli":
		command := "guetzli --quality " + fmt.Sprintf("%d", quality) + ` "` + input + `" "` + output + `"`
		cmd := exec.Command("/bin/sh", "-c", command)
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

// magick cockatoo.tif -clip -negate negated.tif
func (ps *PictureProcessService) Clip(input, output string) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImage(input); err != nil {
		return err
	}

	if err := mw.WriteImage(output); err != nil {
		return err
	}

	return nil
}

func (ps *PictureProcessService) Format(input, output string, format string) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImage(input); err != nil {
		return err
	}

	if err := mw.SetFormat(format); err != nil {
		return err
	}

	if err := mw.WriteImage(output); err != nil {
		return err
	}

	return nil
}

func (ps *PictureProcessService) Exif(input, output string) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.ReadImage(input); err != nil {
		return err
	}

	if err := mw.StripImage(); err != nil {
		return err
	}

	if err := mw.WriteImage(output); err != nil {
		return err
	}

	return nil
}

// convert command
func (ps *PictureProcessService) Convert(input, output string, args ...string) (*imagick.ImageCommandResult, error) {
	imagick.Initialize()
	defer imagick.Terminate()

	commandArr := []string{
		"convert",
		input,
	}
	commandArr = append(append(commandArr, args...), output)
	ret, err := imagick.ConvertImageCommand(commandArr)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
