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