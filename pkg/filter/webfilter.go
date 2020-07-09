package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fuloge/basework/api"
	"github.com/fuloge/basework/configs"
	"github.com/fuloge/basework/pkg/auth"
	"github.com/fuloge/basework/pkg/log"
	"github.com/fuloge/basework/pkg/rsa"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var logger = log.New()

type Filter struct {
}

type loginModul struct {
	Uid int64 `json:"uid"`
}

func (f *Filter) buildResponse(code int, status bool, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":      code,
		"success":   status,
		"data":      data,
		"timestamp": time.Now().Unix(),
	})

	c.Abort()
}

//请求head,必须包含auth,exp项
func (f *Filter) Checkauth() gin.HandlerFunc {
	return func(c *gin.Context) {
		str := fmt.Sprintf("------method:%s\n path:%s", c.Request.Method, c.FullPath())
		fmt.Println(str)

		loginmodul := loginModul{}
		method := c.Request.Method
		switch method {
		case "GET":
			var values = c.Request.URL.Query()
			loginmodul.Uid, _ = strconv.ParseInt(values["uid"][0], 10, 64)
			fmt.Printf("---->parame:%s \n", values)
		case "PUT", "DELETE", "POST":
			data, err := c.GetRawData()
			if configs.EnvConfig.RunMode == 1 {
				fmt.Printf("---->parame:%s \n", data)
			}
			if err != nil {
				logger.Error(api.HTTPParamErr.Message, zap.String(api.HTTPParamErr.Message, err.Error()))
				f.buildResponse(api.HTTPParamErr.Code, false, api.HTTPParamErr.Message, c)
				return
			}

			if err = json.Unmarshal(data, &loginmodul); err != nil {
				logger.Error(api.HTTPParamErr.Message, zap.String(api.HTTPParamErr.Message, err.Error()))
				f.buildResponse(api.HTTPParamErr.Code, false, api.HTTPParamErr.Message, c)
				return
			}

			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		default:
			fmt.Println("no support")
		}

		if _, ok := configs.WhiteList[c.FullPath()]; ok {
			//放行
			c.Next()

			return
		}

		jwt := auth.New()
		a := c.Request.Header.Get("auth")
		if a == "" {
			fmt.Println(api.TokenNilErr.Message)
			f.buildResponse(api.TokenNilErr.Code, false, api.TokenNilErr.Message, c)
			return
		}
		//e := c.Request.Header.Get("exp")

		//isOK = true

		errno := &api.Errno{}
		//解密
		if configs.EnvConfig.RunMode != 1 {
			a, errno = rsa.RsaDecrypt(a)
			if errno != nil {
				f.buildResponse(api.RSADecERR.Code, false, api.RSADecERR.Message, c)
				return
			}
		}

		/*
			expData, err := rsa.RsaDecrypt(e)
			if err != nil {
				f.respondWithError(10002, err.Error(), c)
				return
			}

			//超时
			t, _ := strconv.ParseInt(string(expData), 10, 64)
			if time.Now().Unix() > t {
				f.respondWithError(10003, err.Error(), c)
				return
			}
		*/

		//token
		if !jwt.TokenIsInvalid(a) {
			f.buildResponse(api.TokenInvidErr.Code, false, api.TokenInvidErr.Message, c)
			return
		}

		if m, err := jwt.ParseToken(a); err != nil {
			f.buildResponse(err.Code, false, err.Message, c)
			return
		} else {
			uid := m["uid"].(string)
			Uid, _ := strconv.ParseInt(uid, 10, 64)
			if Uid != loginmodul.Uid {
				f.buildResponse(api.HTTPUidErr.Code, false, api.HTTPUidErr.Message, c)
				return
			}
		}

		//放行
		c.Next()
	}
}
