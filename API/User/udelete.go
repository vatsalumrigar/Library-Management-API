package user

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
)

// @Summary delete one user from user collection
// @ID delete-one-user
// @Produce json
// @Param uId header string true "UserID"
// @Param language header string true "languageToken"
// @Success 201 {object} model.User
// @Failure 500 {object} error
// @Failure 404 {object} error
// @Router /deleteUser/ [delete]
func DeleteUser(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	if middleware.Authentication(c) {

		//userId := c.Param("userId")

		uId , err1 := c.Get("uId")

		if !err1 {
			logs.Error(err1)
			c.JSON(http.StatusNotFound, localization.GetMessage(languageToken,"404"))
			return
		}
		userId := uId.(string)

		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId) 
		result, err := userCollection.DeleteOne(ctx, bson.M{"_Id": objId}) 
		res := map[string]interface{}{"data": result}

		if err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}
		
		if result.DeletedCount < 1 {
			logs.Error("No data to delete")
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"DeleteUser.500"))
			return
		}
		
		c.JSON(http.StatusCreated, gin.H{"message": localization.GetMessage(languageToken,"DeleteUser.201"), "Data": res})
	}

}