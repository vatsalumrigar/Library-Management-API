package middleware

import (
	helper "PR_2/helper"


	"net/http"

	"github.com/gin-gonic/gin"
)

// Authz validates token and authorizes users
func Authentication(c *gin.Context) bool {

        clientToken := c.Request.Header.Get("token")
		

        if clientToken == "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})
            c.Abort()
            return false
        }

        claims, err := helper.ValidateToken(clientToken)

        if err != "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err})
            c.Abort()
            return false 
        }

        c.Set("email", claims.Email)
        c.Set("first_name", claims.Firstname)
        c.Set("last_name", claims.Lastname)
        c.Set("uid", claims.Uid)

        c.Next()

        // fmt.Printf("claims.Uid: %v\n", claims.Uid)
        // fmt.Printf("claims.Firstname: %v\n", claims.Firstname)

    return true

}