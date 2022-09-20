package Middlewares

import (
	"fmt"
	"net/http"
	"test/Helper"

	"github.com/gin-gonic/gin"
)

type Auth struct{}

func (s *Auth) Auth(c *gin.Context) {
	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
		c.Abort()
		return
	}

	_, err := Helper.ValidateToken(clientToken)

	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

}
