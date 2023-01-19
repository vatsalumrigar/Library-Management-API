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
package accounting

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary pay penalty of user
// @ID acc-penalty-pay
// @Accept json
// @Produce json
// @Param librarianId header string true "LibrarianID"
// @Success 201 {object} model.User 
// @Success 201 {object} model.Accounting
// @Failure 400 {object} error
// @Failure 403 {object} error
// @Failure 404 {object} error
// @Failure 406 {object} error
// @Failure 500 {object} error
// @Router /Accounting/penaltypay [post]
func AccountingPenaltyPay(c *gin.Context){

	accountingCollection := database.GetCollection("Accounting")
	userCollection := database.GetCollection("User")
	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	var penaltypay model.PenaltyPay

	if err:= c.BindJSON(&penaltypay); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return

	}

	if middleware.Authentication(c){

		lId,err := c.Get("uid")

		if !err {
			c.JSON(http.StatusNotFound, gin.H{"message": err})
			return
		}

		libId := lId.(string)
		objId1, _ := primitive.ObjectIDFromHex(libId)

		var lib model.User

		err1 := userCollection.FindOne(ctx, bson.M{"_Id": objId1}).Decode(&lib)
		
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err1})
			return
		}

		if lib.UserType != "Librarian"{
			c.JSON(http.StatusForbidden, gin.H{"message": "enter valid librairian token"})
			return
		}

		var user model.User
		
		err2 := userCollection.FindOne(ctx, bson.M{"Username": penaltypay.Username}).Decode(&user)
	
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err2})
			return
		}

		if user.Total_Penalty == 0 {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "user has no pending penalty"})
			return
		}

		if penaltypay.Pay_Amount == user.Total_Penalty {

			match := bson.M{"Username" : penaltypay.Username}
	
			update := bson.M{
					
				"$set": bson.M{
				"Status": "Available",
				"Total_Penalty": 0 ,
				},
	
			}
				
			_, err3 := userCollection.UpdateOne(ctx,match,update)
	
	
				if err3 != nil {
							
					c.JSON(http.StatusInternalServerError, gin.H{"message": err3})
					return
	
				}

			var accountinguser model.Accounting

			err4 := accountingCollection.FindOne(ctx, bson.M{"Email": user.Email}).Decode(&accountinguser)

			if err4 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err4})
				return
			}

			for _, book := range accountinguser.PenaltyDetail {

				booktitle := book.BookTitle

				var foundbook model.Books

				err5 := bookCollection.FindOne(ctx, bson.M{"Title": booktitle}).Decode(&foundbook)
		
				if err5 != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "could not find book"})
					return
				}

				//userbookstaken := user.BooksTaken

				var bd = model.Bookdetails{

							BookId : foundbook.ID.Hex(),
							Title: foundbook.Title,
							TimeTaken: time.Now().Unix(),

				}

				match := bson.M{"Email": user.Email}
				change := bson.M{"$pull": bson.M{"Books_Taken": bson.M{"Title":bd.Title}}}

				_, err := userCollection.UpdateOne(ctx, match, change)
	
				if err != nil {
			
					c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{"error":"cannot update usercollection"})
					return

				}

				qty := foundbook.Quantities
				fmt.Printf("qty: %v\n", qty)

				if book.ReasonType == 1{

					qty +=  1
					fmt.Printf("updtqty: %v\n", qty)
					filter := bson.M{"Title": foundbook.Title}
					update := bson.M {
						"$set": bson.M {
							"Quantities" : qty,
							"Status" : "Available",
						},
					}

		
					_ , err6 := bookCollection.UpdateOne(ctx,filter,update)
		
					if err6!= nil {
		
						c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{"error":"cannot update bookcollection"})
						return

					}

				}//else if book.ReasonType == 2 || book.ReasonType == 3 {
				// 	continue
				// }
			}

			matchaccounting := bson.M{"Email": user.Email}
			updateaccounting := bson.M{
					
				"$set": bson.M{
				"TotalPenalty": 0,
				"TimePenaltyPay": time.Now().Unix() ,
				"PenaltyDetails.$[].PenaltyPay": true,
				},
	
			}

			_, err7 := accountingCollection.UpdateOne(ctx,matchaccounting,updateaccounting)
	
	
				if err7 != nil {
							
					c.JSON(http.StatusInternalServerError, gin.H{"message": err7})
					return
	
				}

	
			c.JSON(http.StatusCreated, gin.H{"message": "penalty payed successfully!"})	
			return
	
		} else {
	
			c.JSON(http.StatusInternalServerError, gin.H{"amount pay should be equal to total penalty": user.Total_Penalty})
			return
	
		}

	}

}