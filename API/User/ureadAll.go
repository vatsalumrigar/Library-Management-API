package user

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	logs "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	localization "PR_2/localise"
)

// @Summary read all user from user collection
// @ID read-all-user
// @Produce json
// @Param queryWord1 query string false "UserType"
// @Param queryWord2 query string false "Firstname"
// @Param language header string true "languageToken"
// @Success 200 {object} model.User
// @Failure 500 {object} error
// @Router /getAllUser/ [get]
func ReadAllUser(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")

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
			logs.Error(err.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		}

		result = append(result, resl)

	}

	c.JSON(http.StatusOK, gin.H{"Data": result})

}

