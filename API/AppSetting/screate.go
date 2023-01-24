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

// @Summary create app setting
// @ID create-setting
// @Accept json
// @Produce json
// @Param adminId header string true "AdminID"
// @Param language header string true "languageToken"
// @Param payload body model.Timings true "Query Payload for create App Timings API"
// @Success 201 {object} model.Timings
// @Failure 400 {object} error
// @Failure 403 {object} error
// @Failure 500 {object} error
// @Router /CreateSetting/ [post]
func CreateSetting(c *gin.Context){	

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
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"CreateSetting.500"))
			return
		}

		if admin.UserType != "Admin" {
			logs.Error("enter valid admin token")
			c.AbortWithStatusJSON(http.StatusForbidden, localization.GetMessage(languageToken,"CreateSetting.403"))
			return
		}
	
		setting := new(model.Timings)

		if err := c.BindJSON(&setting); err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
			return
		}


		addsetting := model.Timings{

			Timing: setting.Timing,

		}

		result, err := appsettingCollection.InsertOne(ctx, addsetting)
		res := map[string]interface{}{"data": result}
	
		if err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}
	
		c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": res})
		return
		

	}

}