package user

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	localization "PR_2/localise"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary read one user from user collection
// @ID read-one-user
// @Produce json
// @Param uId header string true "UserID"
// @Param language header string true "languageToken"
// @Success 200 {object} model.User
// @Failure 500 {object} error
// @Failure 404 {object} error
// @Router /getOneUser/ [get]
func ReadOneUser(c *gin.Context)  {

	languageToken := c.Request.Header.Get("lan")

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	if middleware.Authentication(c) {
	
		//userId := c.Param("userId")

		uId,err1 := c.Get("uid")

		if !err1 {
			logs.Error(err1)
			c.JSON(http.StatusNotFound, localization.GetMessage(languageToken,"404"))
			return
		}

		userId := uId.(string)
		fmt.Printf("userId: %v\n", userId)
		// var result model.User
		var result map[string]interface{}

		objId, _ := primitive.ObjectIDFromHex(userId)

		// opts := options.FindOne().SetProjection(bson.M{
		// 	"_Id":1,
		// 	"User_Type" : 1,
		// 	"Firstname" : 1,
		// 	"Lastname" : 1,
		// 	"Fullname": 1,
		// 	"Email" : 1,
		// 	"Mobile_No" : 1,
		// 	"Password" : 1,
		// 	"Username" : 1,
		// 	"Books_Taken" :1,
		// 	"Status" : 1,
		// 	"Dob" : 1,
		// 	"Login" : 1,
		// 	"IsFirstLogin" : 1,
		// 	"Total_Penalty" : 1,
		// 	"Address" : 1,
		// })
		
		var fullname string

		err := userCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&result)

		if result["Fullname"].(map[string]interface{})[languageToken] == nil {
		
			fullname = result["Fullname"].(map[string]interface{})["en"].(string)

		} else {

			fullname = result["Fullname"].(map[string]interface{})[languageToken].(string)

		}
	
		result["Fullname"] = fullname
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			logs.Error(err.Error())
			return
		}

		//res := map[string]interface{}{"data":result}
		c.JSON(http.StatusOK, gin.H{"message": localization.GetMessage(languageToken,"ReadOneUser.200"), "Data": result})

	}

}
