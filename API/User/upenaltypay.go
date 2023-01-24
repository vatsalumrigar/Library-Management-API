package user

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	localization "PR_2/localise"
	logs "github.com/sirupsen/logrus"
)

// @Summary pay penalty of user
// @ID users-penalty-pay
// @Accept json
// @Produce json
// @Param payload body model.PenaltyPay true "Payload for Penalty Pay API"
// @Param language header string true "languageToken"
// @Success 201 {object} model.User
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /PenaltyPay/ [patch]
func UserPenaltyPay(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	penalty_payer := new(model.PenaltyPay)

	defer cancel()
	
	if err:= c.BindJSON(&penalty_payer); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
		return
	}

	
	var user model.User
	//var book bmodel.Books
	

	err1 := userCollection.FindOne(ctx, bson.M{"Username": penalty_payer.Username}).Decode(&user)

	if err1 != nil {
		logs.Error(err1.Error())
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		return
	}

	if penalty_payer.Pay_Amount == user.Total_Penalty {

		match := bson.M{"Username" : penalty_payer.Username}

		update := bson.M{
				
			"$set": bson.M{
			"Status": "Available",
			"Total_Penalty": 0 ,
			},

		}
			
		_, err := userCollection.UpdateOne(ctx,match,update)


			if err != nil {
					
				logs.Error(err.Error())
				c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
				return

			}

			c.JSON(http.StatusCreated, localization.GetMessage(languageToken,"UserPenaltyPay.201"))	
			return

	} else {

		logs.Error("amount pay should be equal to total penalty:", user.Total_Penalty)
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"UserPenaltyPay.500"))
		return

	}

}