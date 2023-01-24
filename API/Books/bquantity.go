package book

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get book title and their quantities from book collection
// @ID quantity-book
// @Produce json
// @Param language header string true "languageToken"
// @Param librarianId header string true "LibrarianID"
// @Success 200 {object} []map[string]interface{}
// @Failure 403 {string} string 
// @Failure 406 {string} string 
// @Failure 500 {string} string 
// @Router /getQuantityBook/ [get]
func QuantityBook(c *gin.Context){

	languageToken := c.Request.Header.Get("lan")

	userCollection := database.GetCollection("User")
	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	if middleware.Authentication(c){

		lId,err := c.Get("uid")

		if !err {
			logs.Error(err)
			c.JSON(http.StatusNotFound, localization.GetMessage(languageToken,"404"))
			return
		}

		libId := lId.(string)
		objId1, _ := primitive.ObjectIDFromHex(libId)

		var lib model.User

		err1 := userCollection.FindOne(ctx, bson.M{"_Id": objId1}).Decode(&lib)
		
		if err1 != nil {
			logs.Error(err1.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}

		if lib.UserType != "Librarian"{
			logs.Error("enter valid librairian token")
			c.JSON(http.StatusForbidden, localization.GetMessage(languageToken,"QuantityBook.403"))
			return
		}
	
	}

	if !middleware.Authentication(c){

		logs.Error("provide librarian token in header")
		c.AbortWithStatusJSON(http.StatusNotAcceptable, localization.GetMessage(languageToken,"QuantityBook.406"))
		return

	}

	match := bson.M{}
	opts := options.Find().SetProjection(bson.M{"_Id":1,"Title":1, "Quantities":1,"_id":0})

	cursor, _ := bookCollection.Find(ctx, match,opts)
	// var result []model.Books
	var result []map[string]interface{}

	for cursor.Next(ctx){																

		var resl map[string]interface{}
		err := cursor.Decode(&resl)

		if err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}

		result = append(result, resl)
	}

	c.JSON(http.StatusOK, gin.H{"Data": result})

}