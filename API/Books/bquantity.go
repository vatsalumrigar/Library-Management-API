package book

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get book title and their quantities from book collection
// @ID quantity-book
// @Produce json
// @Param librarianId header string true "LibrarianID"
// @Success 200 {object} []map[string]interface{}
// @Failure 403 {string} string 
// @Failure 406 {string} string 
// @Failure 500 {string} string 
// @Router /getQuantityBook/ [get]
func QuantityBook(c *gin.Context){

	userCollection := database.GetCollection("User")
	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	if middleware.Authentication(c){

		lId,err := c.Get("uid")

		if !err {
			logs.Error(err)
			c.JSON(http.StatusNotFound, gin.H{"message": err})
			return
		}

		libId := lId.(string)
		objId1, _ := primitive.ObjectIDFromHex(libId)

		var lib model.User

		err1 := userCollection.FindOne(ctx, bson.M{"_Id": objId1}).Decode(&lib)
		
		if err1 != nil {
			logs.Error(err1.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err1})
			return
		}

		if lib.UserType != "Librarian"{
			logs.Error("enter valid librairian token")
			c.JSON(http.StatusForbidden, gin.H{"message": "enter valid librairian token"})
			return
		}
	
	}

	if !middleware.Authentication(c){
		logs.Error("provide librarian token in header")
		c.AbortWithStatusJSON(http.StatusNotAcceptable,gin.H{"message": "provide librarian token in header"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
			return
		}

		result = append(result, resl)
	}

	c.JSON(http.StatusOK, gin.H{"Data": result})

}