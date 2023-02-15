package accounting

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	localization "PR_2/localise"
	model "PR_2/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary pay penalty of user
// @ID acc-penalty-pay
// @Accept json
// @Produce json
// @Param librarianId header string true "LibrarianID"
// @Param language header string true "languageToken"
// @Param payload body model.PenaltyPay true "Query Payload for Penalty Pay API"
// @Success 201 {object} model.User
// @Success 201 {object} model.Accounting
// @Failure 400 {object} error
// @Failure 403 {object} error
// @Failure 404 {object} error
// @Failure 406 {object} error
// @Failure 500 {object} error
// @Router /Accounting/penaltypay [post]
func AccountingPenaltyPay(c *gin.Context){

	languageToken := c.Request.Header.Get("lan")

	accountingCollection := database.GetCollection("Accounting")
	userCollection := database.GetCollection("User")
	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	var penaltypay model.PenaltyPay

	if err:= c.BindJSON(&penaltypay); err != nil {

		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
		return

	}

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

		fmt.Printf("lib.Firstname: %v\n", lib.Firstname)
		
		if err1 != nil {
			logs.Error(err1.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}

		if lib.UserType != "Librarian"{
			logs.Error("enter valid librairian token")
			c.JSON(http.StatusForbidden, localization.GetMessage(languageToken,"AccountingPenaltyPay.403"))
			return
		}

		var user model.User
		
		err2 := userCollection.FindOne(ctx, bson.M{"Username": penaltypay.Username}).Decode(&user)
	
		if err2 != nil {
			logs.Error(err2.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}

		fmt.Printf("user.Firstname: %v\n", user.Firstname)

		if user.Total_Penalty == 0 {
			logs.Error("user has no pending penalty")
			c.AbortWithStatusJSON(http.StatusNotAcceptable, localization.GetMessage(languageToken,"AccountingPenaltyPay.406"))
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
						
					logs.Error(err3.Error())
					c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
					return
	
				}

			var accountinguser model.Accounting

			err4 := accountingCollection.FindOne(ctx, bson.M{"Email": user.Email}).Decode(&accountinguser)

			if err4 != nil {
				logs.Error(err4.Error())
				c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
				return
			}

			for _, book := range accountinguser.PenaltyDetail {

				booktitle := book.BookTitle

				var foundbook model.Books

				err5 := bookCollection.FindOne(ctx, bson.M{"Title": booktitle}).Decode(&foundbook)
		
				if err5 != nil {
					logs.Error(err5.Error())
					c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"AccountingPenaltyPay.500.error1"))
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
					logs.Error(err.Error())
					c.AbortWithStatusJSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"AccountingPenaltyPay.500.error2"))
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
		
						logs.Error(err6.Error())
						c.AbortWithStatusJSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"AccountingPenaltyPay.500.error3"))
						return

					}

				}

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
							
					logs.Error(err7.Error())
					c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
					return
	
				}

	
			c.JSON(http.StatusCreated, localization.GetMessage(languageToken,"AccountingPenaltyPay.201"))	
			return
	
		} else {
	
			logs.Error("amount pay should be equal to total penalty:", user.Total_Penalty)
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"AccountingPenaltyPay.500.error4"))
			return
	
		}

	}

}