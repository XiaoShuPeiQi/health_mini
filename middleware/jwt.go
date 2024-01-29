package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"main.go/model/common/response"
	"main.go/service"
	"strings"
	"time"
)

var manageAdminUserTokenService = service.ServiceGroupApp.ManageServiceGroup.ManageAdminUserTokenService
var mallUserTokenService = service.ServiceGroupApp.MallServiceGroup.MallUserTokenService

func AdminJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//token := c.Request.Header.Get("token")
		//token := c.GetHeader("Authorization")
		//
		//fmt.Println(token + "..................")
		//if token == "" {
		//	response.FailWithDetailed(nil, "未登录或非法访问", c)
		//	c.Abort()
		//	return
		//}
		token, err := getToken(c)
		if err != nil {
			response.FailWithDetailed(nil, "token格式错误！", c)
			c.Abort()
			return
		}

		err, mallAdminUserToken := manageAdminUserTokenService.ExistAdminToken(token)
		if err != nil {
			response.FailWithDetailed(nil, "未登录或非法访问", c)
			c.Abort()
			return
		}
		if time.Now().After(mallAdminUserToken.ExpireTime) {
			response.FailWithDetailed(nil, "授权已过期", c)
			err = manageAdminUserTokenService.DeleteMallAdminUserToken(token)
			if err != nil {
				return
			}
			c.Abort()
			return
		}
		c.Next()
	}

}

func UserJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//token := c.Request.Header.Get("token")
		//if token == "" {
		//	response.UnLogin(nil, c)
		//	c.Abort()
		//	return
		//}
		token, err2 := getToken(c)
		if err2 != nil {
			response.FailWithDetailed(nil, "token格式错误！", c)
			c.Abort()
			return
		}
		err, mallUserToken := mallUserTokenService.ExistUserToken(token)
		if err != nil {
			response.UnLogin(nil, c)
			c.Abort()
			return
		}
		if time.Now().After(mallUserToken.ExpireTime) {
			response.FailWithDetailed(nil, "授权已过期", c)
			err = mallUserTokenService.DeleteMallUserToken(token)
			if err != nil {
				return
			}
			c.Abort()
			return
		}
		c.Next()
	}

}

// getToken 处理token到想要的值
func getToken(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("Invalid Authorization header format")
	}
	fmt.Println("userToken:" + tokenParts[1])
	fmt.Println("----------------------")
	return tokenParts[1], nil
}
