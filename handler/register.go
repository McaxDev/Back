package handler

import (
	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

func Register(c *gin.Context) {
	userName, password := c.PostForm("userName"), c.PostForm("password")
	if err := passwordvalidator.Validate(password, 60.0); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "password not safe"})
	}
}
