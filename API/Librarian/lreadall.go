package librarians

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	localization "PR_2/localise"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary read all librarian
// @ID read-all-librarian
// @Produce json
// @Param language header string true "languageToken"
// @Success 200 {object} model.User
// @Failure 500 {object} error
// @Router /getAllLibrarian/ [get]
func ReadAllLibrarian(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")

	librarianCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	match :=  bson.M{"User_Type": "Librarian"}
	// opts := options.Find().SetSort(bson.D{{Key: "Firstname",Value: 1}}) // sort
	// opts := options.Find().SetLimit(2) //limit
	// opts :=  options.Find().SetSkip(1) //skip
	// opts := options.Find().SetProjection(bson.M{"_Id": 1,"Email": 1,"Password": 1}) // include field projection
	opts := options.Find().SetProjection(bson.D{{Key: "_Id",Value: 0},{Key: "Email",Value: 0},{Key: "Password",Value: 0}}) // exclude field projection
	
	
	cursor, _ := librarianCollection.Find(ctx,match,opts)
	var result []model.Librarian

	for cursor.Next(ctx){

		var resl model.Librarian
		err := cursor.Decode(&resl)

		if err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}

		result = append(result, resl)

	}

	c.JSON(http.StatusOK, gin.H{"Data": result})

}
