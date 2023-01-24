package librarians

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
)

// @Summary read one librarian
// @ID read-one-librarian
// @Produce json
// @Param language header string true "languageToken"
// @Param librarianId path string true "LibrarianID" 
// @Success 200 {object} model.User
// @Failure 500 {object} error
// @Router /getOneLibrarian/{librarianId} [get]
func ReadOneLibrarian(c *gin.Context)  {

	languageToken := c.Request.Header.Get("lan")

	librarianCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	librarianId := c.Param("librarianId")
	var result model.Librarian


	objId, _ := primitive.ObjectIDFromHex(librarianId)

	err := librarianCollection.FindOne(ctx, bson.M{"_Id": objId,"UserType": "Librarian"}).Decode(&result)
	

	if err != nil {
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		logs.Error(err.Error())
		return
	}
	//res := map[string]interface{}{"data":result}
	c.JSON(http.StatusOK, gin.H{"message": localization.GetMessage(languageToken,"ReadOneLibrarian.200"), "Data": result})

}
