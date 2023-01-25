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
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
)

// @Summary update user in user collection
// @ID update-user
// @Accept json
// @Produce json
// @Param uId header string true "UserID"
// @Param payload body model.User true "Payload for update user API"
// @Param language header string true "languageToken"
// @Success 201 {object} model.User
// @Failure 500 {object} error
// @Failure 409 {object} error
// @Router /updateUser/ [put]
func UpdateUser(c *gin.Context) {
	
	languageToken := c.Request.Header.Get("lan")

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	
	if middleware.Authentication(c) {

		uId,err3 := c.Get("uid")

		if !err3 {
			c.JSON(http.StatusNotFound, localization.GetMessage(languageToken,"404"))
			logs.Error(err3)
			return
		}

		userId := uId.(string)

	
		var user model.User
		
		defer cancel()
		
		objId, _ := primitive.ObjectIDFromHex(userId)
		
		if err := c.BindJSON(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
			logs.Error(err.Error())
			return
		}

		hashpwd := controllers.HashPassword(user.Password)

		edited := bson.M {
			"UserType" : user.UserType,
			"Firstname" : user.Firstname,
			"Lastname" : user.Lastname,
			"Fullname": user.Fullname,
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
			c.AbortWithStatusJSON(http.StatusConflict, localization.GetMessage(languageToken,"409"))
			logs.Error(val.Error())
			return
		} 

		result, err := userCollection.UpdateOne(ctx, bson.M{"_Id": objId}, bson.M{"$set": edited})

		res := map[string]interface{}{"data": result}

		if err != nil {
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			logs.Error(err.Error())
			return
			}

		if result.MatchedCount < 1 {
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"UpdateUser.500"))
			logs.Error("Data doesn't exist")
			return
		}
				
		c.JSON(http.StatusCreated, gin.H{"message": localization.GetMessage(languageToken,"UpdateUser.201"), "Data": res})	
	}
}

