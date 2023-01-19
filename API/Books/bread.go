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
package book

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary read book from book collection
// @ID read-book
// @Produce json
// @Param bookId path string true "BookID"
// @Success 200 {object} model.Books 
// @Failure 500 {string} string 
// @Router /getOneBook/{bookId} [get]
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

