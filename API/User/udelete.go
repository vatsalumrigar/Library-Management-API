package user

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary delete one user from user collection
// @ID delete-one-user
// @Produce json
// @Param uId header string true "UserID"
// @Success 201 {object} model.User
// @Failure 500 {object} error
// @Failure 404 {object} error
// @Router /deleteUser/ [delete]
func DeleteUser(c *gin.Context) {

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	if middleware.Authentication(c) {

		//userId := c.Param("userId")

		uId , err1 := c.Get("uId")

		if !err1 {
			c.JSON(http.StatusNotFound, gin.H{"message": err1})
			return
		}
		userId := uId.(string)

		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId) 
		result, err := userCollection.DeleteOne(ctx, bson.M{"_Id": objId}) 
		res := map[string]interface{}{"data": result}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		
		if result.DeletedCount < 1 {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "No data to delete"})
			return
		}
		
		c.JSON(http.StatusCreated, gin.H{"message": "Article deleted successfully", "Data": res})
	}

}