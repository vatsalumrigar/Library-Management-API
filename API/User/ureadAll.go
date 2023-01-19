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
package user

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary read all user from user collection
// @ID read-all-user
// @Produce json
// @Param queryWord1 query string false "UserType"
// @Param queryWord2 query string false "Firstname"
// @Success 200 {object} model.User
// @Failure 500 {object} error
// @Router /getAllUser/ [get]
func ReadAllUser(c *gin.Context) {

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	queryWord := c.Request.URL.Query().Get("queryWord")

	match :=  bson.M{"User_Type":"User"}

	if queryWord != "" {

		match = bson.M{
			"User_Type":"User",
			"Firstname": bson.M{"$regex": queryWord, "$options": "im"},
		}
	}

	// opts := options.Find().SetSort(bson.D{{Key: "Firstname",Value: 1}}) // sort
	// opts := options.Find().SetLimit(2) //limit
	// opts :=  options.Find().SetSkip(1) //skip
	// opts := options.Find().SetProjection(bson.M{"_Id": 1,"Email": 1,"Books_Taken": 1}) // include field projection
	// opts := options.Find().SetProjection(bson.D{{Key: "_Id",Value: 0},{Key: "Email",Value: 0}}) // exclude field projection
	
	cursor, _ := userCollection.Find(ctx,match)
	var result []model.User

	for cursor.Next(ctx){

		var resl model.User
		err := cursor.Decode(&resl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message":err})
		}

		result = append(result, resl)

	}

	c.JSON(http.StatusOK, gin.H{"Data": result})

}

