package admin

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ReadAllAdmin(c *gin.Context) {

	adminCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	match :=  bson.M{"User_Type": "Admin"}
	opts := options.Find().SetProjection(bson.D{{Key: "_Id",Value: 0},{Key: "Email",Value: 0},{Key: "Password",Value: 0}}) // exclude field projection
	
	cursor, _ := adminCollection.Find(ctx,match,opts)
	var result []model.Admin

	for cursor.Next(ctx){

		var resl model.Admin
		err := cursor.Decode(&resl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message":err})
		}

		result = append(result, resl)

	}

	c.JSON(http.StatusOK, gin.H{"Data": result})

}
