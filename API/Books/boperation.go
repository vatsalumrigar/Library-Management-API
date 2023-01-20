package book

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OperationBook(c *gin.Context) {

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
	
	var books model.Bookqty


	if err := c.BindJSON(&books); err != nil {

		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return

	}

	if books.Operations != "Add" && books.Operations != "Remove" {

		c.AbortWithStatusJSON(http.StatusNotAcceptable,gin.H{"message":"book opertions should either be: Add or Remove"})
		return

	}

	for title , qty := range books.Books{

		fmt.Printf("title: %v\n", title)
		fmt.Printf("book: %v\n", qty)

		titleCount, _ := bookCollection.CountDocuments(ctx, bson.M{"Title": title})

		if titleCount < 1 {
			logs.Error(title,"-no such book")
			c.AbortWithStatusJSON(http.StatusNotFound,gin.H{"error": title+"-no such book"})
			return
		}
	
	}


}