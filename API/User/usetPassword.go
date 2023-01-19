// @title Library Management API
// @version 1.0
// @description This is a  Library Management API server.
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /
// @query.collection.format multi
package user

import (
	controllers "PR_2/Controller"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary set new password for user 
// @ID user-set-new-password
// @Accept json
// @Produce json
// @Success 201 {object} model.User
// @Failure 400 {object} error
// @Failure 406 {object} error
// @Failure 500 {object} error
// @Router /UserSetNewPassword/ [patch]
func SetNewPasswordUser(c *gin.Context){

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	userpwd := new(model.SetNewPassword)

	defer cancel()

	if err:= c.BindJSON(&userpwd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user model.User

	err := userCollection.FindOne(ctx, bson.M{"Email": userpwd.Email}).Decode(&user)		

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if user.Password != userpwd.OldPassword {
		c.AbortWithStatusJSON(http.StatusNotAcceptable,gin.H{"error":"old password is incorrect"})
		return
	}

	hashpwd := controllers.HashPassword(userpwd.NewPassword)

	match := bson.M{"Email": userpwd.Email}
	update := bson.M{"Password": hashpwd, "IsFirstLogin": false}

	result, err := userCollection.UpdateOne(ctx,match,bson.M{"$set":update})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	res := map[string]interface{}{"data":result}

	if result.ModifiedCount < 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Data doesn't exist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "data updated successfully!", "Data": res})	

}