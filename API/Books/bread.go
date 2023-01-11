package book

import (
	model "PR_2/model"
	database "PR_2/databases"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReadOneBook(c *gin.Context) {

	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)
	bookId := c.Param("bookId")
	var result model.Books
	
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bookId)

	err := bookCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&result)
	res := map[string]interface{}{"data":result}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Fetched!", "Data": res})

}

