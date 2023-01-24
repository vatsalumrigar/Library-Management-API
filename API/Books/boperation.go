package book

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary update book quantites according to operation in book collection
// @ID operation-book
// @Accept json
// @Produce json
// @Param language header string true "languageToken"
// @Param librarianId header string true "LibrarianID"
// @Success 201 {object} string
// @Failure 404 {string} string 
// @Failure 406 {string} string 
// @Failure 500 {string} string 
// @Router /operationBook/ [patch]
func OperationBook(c *gin.Context) {

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
			c.JSON(http.StatusNotAcceptable, localization.GetMessage(languageToken,"OperationBook.406.error1"))
			return
		}
	
	}

	if !middleware.Authentication(c){
		logs.Error("provide librarian token in header")
		c.AbortWithStatusJSON(http.StatusNotAcceptable, localization.GetMessage(languageToken,"OperationBook.406.error2"))
		return
	}
	
	var books model.Bookqty


	if err := c.BindJSON(&books); err != nil {

		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		return

	}

	if books.Operations != "Add" && books.Operations != "Subtract" {

		c.AbortWithStatusJSON(http.StatusNotAcceptable,localization.GetMessage(languageToken,"OperationBook.406.error3"))
		return

	}

	for title , qty := range books.Books{

		// fmt.Printf("title: %v\n", title)
		// fmt.Printf("qty: %v\n", qty)

		operationQty := qty.(string)
 		operationQuantity, _ := strconv.Atoi(operationQty)


		var foundBook model.Books

		match := bson.M{"Title": title}
		update := bson.M{}

		err := bookCollection.FindOne(ctx, match).Decode(&foundBook)
	
		if err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}

		if books.Operations == "Add" {

			update = bson.M{
					
				"$set": bson.M{
				"Status": "Available",
				"Quantities": foundBook.Quantities + operationQuantity ,

				},
	
			}

		}

		if books.Operations == "Subtract" {

			if foundBook.Quantities - operationQuantity == 0 {

				update = bson.M{
					
					"$set": bson.M{
					"Status": "Unavailable",
					"Quantities": 0 ,
	
					},
		
				}

			}

			if foundBook.Quantities - operationQuantity < 0 {

				logs.Error("books quantity cannot be less then zero for book:", title)
				c.AbortWithStatusJSON(http.StatusNotAcceptable, localization.GetMessage(languageToken,"OperationBook.406.error4"))
				return

			}

			update = bson.M{
					
				"$set": bson.M{

				"Quantities": foundBook.Quantities - operationQuantity ,

				},
	
			}

		}
			
		_, err = bookCollection.UpdateOne(ctx,match,update)


			if err != nil {
					
				logs.Error(err.Error())
				c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
				return

			}
	
	}

	c.JSON(http.StatusCreated,gin.H{"message":localization.GetMessage(languageToken,"OperationBook.201")})

}