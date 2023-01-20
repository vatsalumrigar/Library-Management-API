package user

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary pay penalty of user
// @ID users-penalty-pay
// @Accept json
// @Produce json
// @Param payload body model.PenaltyPay true "Payload for Penalty Pay API"
// @Success 201 {object} model.User
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /PenaltyPay/ [patch]
func UserPenaltyPay(c *gin.Context) {

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	penalty_payer := new(model.PenaltyPay)

	defer cancel()
	
	if err:= c.BindJSON(&penalty_payer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	
	var user model.User
	//var book bmodel.Books
	

	err1 := userCollection.FindOne(ctx, bson.M{"Username": penalty_payer.Username}).Decode(&user)

	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err1})
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
						
				c.JSON(http.StatusInternalServerError, gin.H{"message": err})
				return

			}

			c.JSON(http.StatusCreated, gin.H{"message": "penalty payed and book returned successfully successfully!"})	
			return

	} else {

		c.JSON(http.StatusInternalServerError, gin.H{"amount pay should be equal to total penalty": user.Total_Penalty})
		return

	}

}