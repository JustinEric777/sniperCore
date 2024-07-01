package untils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

/**
 * 发送http请求
 */
func ReqClient(method string, url string, headers map[string]string, params interface{}) (string, error) {
	var request *http.Request
	var err error
	client := &http.Client{}

	byteData, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	if strings.ToLower(method) == "get" {
		request, err = http.NewRequest("GET", url, nil)
	} else if strings.ToLower(method) == "post" {
		request, err = http.NewRequest("POST", url, bytes.NewBuffer(byteData))
	} else {
		err = errors.New("only support get and post request")
	}
	if err != nil {
		return "", err
	}

	if len(headers) != 0 {
		for key, val := range headers {
			request.Header.Add(key, val)
		}
	}

	resp, _ := client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}
