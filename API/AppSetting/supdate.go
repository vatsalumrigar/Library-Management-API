package appsetting

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary update app setting
// @ID update-setting
// @Accept json
// @Produce json
// @Param adminId header string true "AdminID"
// @Param language header string true "languageToken"
// @Param payload body model.Timings true "Query Payload for update App Timings API"
// @Success 201 {object} model.Timings
// @Failure 400 {object} error
// @Failure 403 {object} error
// @Failure 404 {object} error
// @Failure 406 {object} error
// @Failure 500 {object} error
// @Router /UpdateSetting/ [put]
func UpdateSetting(c *gin.Context){

	languageToken := c.Request.Header.Get("lan")

	if middleware.Authentication(c){

		appsettingCollection := database.GetCollection("AppSetting")
		userCollection := database.GetCollection("User")
		ctx, cancel := database.DbContext(10)

		defer cancel()

		aId,_ := c.Get("uid")

		var admin model.User

		adminId := aId.(string)	
		objId, _ := primitive.ObjectIDFromHex(adminId)

		err := userCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&admin)

		if err != nil {

			logs.Error(err.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"UpdateSetting.500.error1"))
			return

		}

		if admin.UserType != "Admin" {

			logs.Error("enter valid admin token")
			c.AbortWithStatusJSON(http.StatusForbidden, localization.GetMessage(languageToken,"UpdateSetting.403"))
			return

		}

		setting := new(model.Timings)

		if err := c.BindJSON(&setting); err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
			return
		}

		edited := bson.M{
			"Timing": setting.Timing,
		}

		id, _ := primitive.ObjectIDFromHex("63c1398db715e13c41025cb8")

		result, err := appsettingCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": edited})

		res := map[string]interface{}{"data": result}

		if err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
			}

		if result.MatchedCount < 1 {
			logs.Error("Data doesn't exist")
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"UpdateSetting.500.error2"))
			return
		}
				
		c.JSON(http.StatusCreated, gin.H{"message": "data updated successfully!", "Data": res})	

		
	}

}