package admin

import (
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReadOneAdmin(c *gin.Context)  {

	adminCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	adminId := c.Param("adminId")
	var result model.Admin


	objId, _ := primitive.ObjectIDFromHex(adminId)

	err := adminCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&result)
	

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		fmt.Println(err)
		return
	}
	//res := map[string]interface{}{"data":result}
	c.JSON(http.StatusOK, gin.H{"message": "Data Fetched!", "Data": result})

}
