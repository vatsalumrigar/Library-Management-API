package book

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary get issued books
// @ID book-issued
// @Produce json
// @Param uId header string true "UserID"
// @Param language header string true "languageToken"
// @Success 200 {object}  model.BooksIssued
// @Failure 403 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /getIssuedBook/ [get]
func IssuedBook(c *gin.Context){

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
			c.JSON(http.StatusForbidden, localization.GetMessage(languageToken,"IssuedBook.403"))
			return
		}

		cursor,_ := bookCollection.Find(ctx, bson.M{})
		var result []model.BooksIssued

		for cursor.Next(ctx) {

			var book model.Books

			err:= cursor.Decode(&book)
			if err != nil{

				logs.Error(err.Error())
				c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
				return

			}

			fmt.Printf("book.Title: %v\n", book.Title)

			var issueto []model.IssueDetails
			var totalissueb int

			find, _ :=  userCollection.Find(ctx, bson.M{"Books_Taken.Title": book.Title})

			for find.Next(ctx) {

				var user model.User

				er:= find.Decode(&user)

				if er != nil {
	
					logs.Error(er.Error())
					c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
					return
	
				}

				fmt.Printf("user.Email: %v\n", user.Email)
				qty := 0
				for _,b := range user.BooksTaken{
					if b.Title == book.Title {
						qty += 1
					}
				}

				issuedetails := model.IssueDetails {
					UserID: user.ID.Hex(),
					Email: user.Email,
					Quantity: qty ,
				}

				issueto = append(issueto, issuedetails )

				totalissueb = 0
				for _, issuedetail := range issueto {
					totalissueb += issuedetail.Quantity
				}

			}


			if len(issueto) == 0 {

				continue

			}else {

				bissued := model.BooksIssued {
						BookID: book.ID.Hex(),
						BookTitle: book.Title,
						Cost: book.Cost,
						IssuedTo: issueto ,
						IssuedQuantity: totalissueb ,
						BooksLeft: book.Quantities,
					}

				result = append(result, bissued)

			}

		}

		c.JSON(http.StatusOK, gin.H{"Data": result})
		return
	
	}

}