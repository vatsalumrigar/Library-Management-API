package book

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary update book quantites according to operation in book collection
// @ID operation-book
// @Accept json
// @Produce json
// @Param librarianId header string true "LibrarianID"
// @Success 201 {object} string
// @Failure 404 {string} string 
// @Failure 406 {string} string 
// @Failure 500 {string} string 
// @Router /operationBook/ [patch]
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
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "enter valid librairian token"})
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

	if books.Operations != "Add" && books.Operations != "Subtract" {

		c.AbortWithStatusJSON(http.StatusNotAcceptable,gin.H{"message":"book opertions should either be: Add or Subtract"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
				c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error":"books quantity cannot be less then zero for book:"+title})
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
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return

			}
	
	}

	c.JSON(http.StatusCreated,gin.H{"message":"books updated succesfully!"})

}