package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (conf *Configuration) startHttpListener() {
	r := gin.Default()
	r.GET("/api/user", findAllUsers)
	r.GET("/api/user/:id", findUserById)
	r.POST("/api/user", addUser)
	r.PUT("/api/user", updateUser)
	r.PUT("/api/user/:d", updateUser)
	r.DELETE("/api/user/:id", deleteUser)

	r.Run(fmt.Sprintf(":%04d", conf.Http.Port))
}

type Msg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func renderError(c *gin.Context, code int, msg string) {
	m := Msg{
		Code:    code,
		Message: msg,
	}
	c.JSON(code, m)
}
