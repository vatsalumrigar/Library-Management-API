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

func UpdateSetting(c *gin.Context){

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

		edited := bson.M{
			"Timing": setting.Timing,
		}

		id, _ := primitive.ObjectIDFromHex("63c1398db715e13c41025cb8")

		result, err := appsettingCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": edited})

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