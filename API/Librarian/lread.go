package librarians

import (
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReadOneLibrarian(c *gin.Context)  {

	librarianCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	librarianId := c.Param("librarianId")
	var result model.Librarian


	objId, _ := primitive.ObjectIDFromHex(librarianId)

	err := librarianCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&result)
	

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		fmt.Println(err)
		return
	}
	//res := map[string]interface{}{"data":result}
	c.JSON(http.StatusOK, gin.H{"message": "Data Fetched!", "Data": result})

}