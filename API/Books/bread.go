package book

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary read book from book collection
// @ID read-book
// @Produce json
// @Param language header string true "languageToken"
// @Param bookId path string true "BookID"
// @Success 200 {object} model.Books 
// @Failure 500 {string} string 
// @Router /getOneBook/{bookId} [get]
func ReadOneBook(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")

	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)
	bookId := c.Param("bookId")
	var result model.Books
	
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bookId)

	err := bookCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&result)
	res := map[string]interface{}{"data":result}

	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": localization.GetMessage(languageToken,"200"), "Data": res})

}

