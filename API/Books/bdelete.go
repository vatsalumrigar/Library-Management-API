package book

import (
	database "PR_2/databases"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
)

// @Summary delete book from book collection 
// @ID delete-book
// @Produce json
// @Param bookId path string true "BookID" 
// @Success 201 {object} model.Books
// @Failure 500 {object} error
// @Router /deleteBook/{bookId} [delete]
func DeleteBook(c *gin.Context) {

	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)

	bookId := c.Param("bookId")

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bookId) 
	result, err := bookCollection.DeleteOne(ctx, bson.M{"_Id": objId}) 
	res := map[string]interface{}{"data": result}

	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	
	if result.DeletedCount < 1 {
		logs.Error("No data to delete")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No data to delete"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Article deleted successfully", "Data": res})

}

