package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sniperCore/core/log"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func LogToFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		method := ctx.Request.Method
		url := ctx.Request.URL
		clientIp := ctx.ClientIP()
		status := ctx.Writer.Status()
		userAgent := ctx.Request.UserAgent()
		httpVersion := ctx.Request.Proto
		params, _ := getRequestParams(ctx)
		paramsInfo, _ := json.Marshal(params)
		data, _ := ctx.GetRawData()
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		ctx.Next()

		endTime := time.Now()
		cost := endTime.Sub(startTime)

		logInfoStr := fmt.Sprintf(`- %s %s %s %s %s %d %s "%s" %s`,
			clientIp,
			httpVersion,
			method,
			url,
			string(paramsInfo),
			status,
			cost,
			userAgent,
			blw.body.String(),
		)
		if status == http.StatusOK {
			log.Info(logInfoStr)
		} else if status == http.StatusInternalServerError || status == http.StatusBadRequest {
			log.Error(logInfoStr)
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func getRequestParams(c *gin.Context) (map[string]interface{}, error) {

	const defaultMemory = 32 << 20
	contentType := c.ContentType()

	var (
		dataMap  = make(map[string]interface{})
		queryMap = make(map[string]interface{})
		postMap  = make(map[string]interface{})
	)

	for k := range c.Request.URL.Query() {
		queryMap[k] = c.Query(k)
	}

	if "application/json" == contentType {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		if c.Request != nil && c.Request.Body != nil {
			if err := json.NewDecoder(c.Request.Body).Decode(&postMap); err != nil {
				return nil, err
			}
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	} else if "multipart/form-data" == contentType {
		if err := c.Request.ParseMultipartForm(defaultMemory); err != nil {
			return nil, err
		}
		for k, v := range c.Request.PostForm {
			if len(v) > 1 {
				postMap[k] = v
			} else if len(v) == 1 {
				postMap[k] = v[0]
			}
		}
	} else {
		if err := c.Request.ParseForm(); err != nil {
			return nil, err
		}
		if err := c.Request.ParseMultipartForm(defaultMemory); err != nil {
			if err != http.ErrNotMultipart {
				return nil, err
			}
		}
		for k, v := range c.Request.PostForm {
			if len(v) > 1 {
				postMap[k] = v
			} else if len(v) == 1 {
				postMap[k] = v[0]
			}
		}
	}

	var mu sync.RWMutex
	for k, v := range queryMap {
		mu.Lock()
		dataMap[k] = v
		mu.Unlock()
	}
	for k, v := range postMap {
		mu.Lock()
		dataMap[k] = v
		mu.Unlock()
	}

	return dataMap, nil
}
