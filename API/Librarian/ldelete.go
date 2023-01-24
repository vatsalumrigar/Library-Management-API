package librarians

import (
	database "PR_2/databases"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
)

// @Summary delete librarian
// @ID delete-librarian
// @Produce json
// @Param language header string true "languageToken"
// @Param librarianId path string true "LibrarianID" 
// @Success 201 {object} model.User
// @Failure 500 {object} error
// @Router /deleteLibrarian/{librarianId} [delete]
func DeleteLibrarian(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")

	librarianCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	librarianId := c.Param("librarianId")

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(librarianId) 
	result, err := librarianCollection.DeleteOne(ctx, bson.M{"_Id": objId}) 
	res := map[string]interface{}{"data": result}

	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		return
	}
	
	if result.DeletedCount < 1 {
		logs.Error("No data to delete")
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"DeleteLibrarian.500"))
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": localization.GetMessage(languageToken,"DeleteLibrarian.201"), "Data": res})

}