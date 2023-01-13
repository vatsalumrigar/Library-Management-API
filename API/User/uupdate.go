package user

import (
	controllers "PR_2/Controller"
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	validation "PR_2/validation"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateUser(c *gin.Context) {
	
	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	
	if middleware.Authentication(c) {

		uId,err3 := c.Get("uid")

		if !err3 {
			c.JSON(http.StatusNotFound, gin.H{"message": err3})
			return
		}

		userId := uId.(string)

	
		var user model.User
		
		defer cancel()
		
		objId, _ := primitive.ObjectIDFromHex(userId)
		
		if err := c.BindJSON(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		hashpwd := controllers.HashPassword(user.Password)

		edited := bson.M {
			"UserType" : user.UserType,
			"Firstname" : user.Firstname,
			"Lastname" : user.Lastname,
			"Email" : user.Email ,
			"MobileNo" : user.MobileNo ,
			"Password" : hashpwd ,
			"Username" : user.Username ,
			"BooksTaken" : user.BooksTaken ,
			"Status" : user.Status,
			"Dob" : user.Dob,
			"Login": user.Login,
			"Total_Penalty": user.Total_Penalty,
		}
		
		val := validation.ValidateUmodel(ctx, user.Email, user.Username, user.MobileNo, user.Dob, user.Status)

		if val != nil {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": val.Error() })
			return
		} 

		result, err := userCollection.UpdateOne(ctx, bson.M{"_Id": objId}, bson.M{"$set": edited})

		res := map[string]interface{}{"data": result}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
			}

		if result.MatchedCount < 1 {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Data doesn't exist"})
			return
		}
				
		c.JSON(http.StatusCreated, gin.H{"message": "data updated successfully!", "Data": res})	
	}
}
