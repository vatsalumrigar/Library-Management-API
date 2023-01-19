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
package appsetting

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary create app setting
// @ID create-setting
// @Accept json
// @Produce json
// @Param adminId header string true "AdminID"
// @Success 201 {object} model.Timings
// @Failure 400 {object} error
// @Failure 403 {object} error
// @Failure 500 {object} error
// @Router /CreateSetting/ [post]
func CreateSetting(c *gin.Context){	

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
			c.JSON(http.StatusInternalServerError, gin.H{"message": "user not logged in"})
			return
		}

		if admin.UserType != "Admin" {
			c.AbortWithStatusJSON(http.StatusForbidden,gin.H{"error": "enter valid admin token"})
			return
		}
	
		setting := new(model.Timings)

		if err := c.BindJSON(&setting); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}


		addsetting := model.Timings{

			Timing: setting.Timing,

		}

		result, err := appsettingCollection.InsertOne(ctx, addsetting)
		res := map[string]interface{}{"data": result}
	
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error() })
			return
		}
	
		c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": res})
		return
		

	}

}