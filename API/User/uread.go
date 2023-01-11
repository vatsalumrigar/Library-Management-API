package user

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReadOneUser(c *gin.Context)  {

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	if middleware.Authentication(c) {
	
		//userId := c.Param("userId")

		uId,err1 := c.Get("uid")

		if !err1 {
			c.JSON(http.StatusNotFound, gin.H{"message": err1})
			return
		}

		userId := uId.(string)
		fmt.Printf("userId: %v\n", userId)
		var result model.User

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := userCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&result)
		

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			fmt.Println(err)
			return
		}
		//res := map[string]interface{}{"data":result}
		c.JSON(http.StatusOK, gin.H{"message": "Data Fetched!", "Data": result})

	}

}
