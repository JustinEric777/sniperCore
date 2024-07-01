package untils

import (
	"fmt"
	"math"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

/**
 * 判断是否post请求
 */
func IsPost(ctx *gin.Context) bool {
	if method := ctx.Request.Method; method == "POST" {
		return true
	}
	return false
}

func IsSuperUser(roles string) bool {
	roleArr := strings.Split(roles, ",")
	for _, role := range roleArr {
		if role == "super" {
			return true
		}
	}

	return false
}

func ShellExecute(command string) (string, error) {
	fmt.Println(command)
	cmd := exec.Command("/bin/bash", "-c", command)

	output, err := cmd.Output()
	if err != nil {
		return string(output), fmt.Errorf("Execute Shell:%s failed with error:%s", command, err.Error())
	}

	return string(output), nil
}

func VoiceCallParams(text string) []string {
	var templateParams []string
	count := int(math.Ceil(float64(len(text)) / 32))
	for i := 0; i < 4; i++ {
		if i < count {
			var tempText string
			if len(text[i*32:]) >= 32 {
				tempText = text[i*32 : (i+1)*32]
			} else {
				tempText = text[i*32:]
			}
			templateParams = append(templateParams, tempText)
		} else {
			templateParams = append(templateParams, "")
		}
	}

	return templateParams
}
