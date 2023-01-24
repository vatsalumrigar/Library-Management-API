package admin

import (
	database "PR_2/databases"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
)

// @Summary delete admin
// @ID delete-admin
// @Produce json
// @Param language header string true "languageToken"
// @Param adminId path string true "AdminID" 
// @Success 201 {object} model.User
// @Failure 500 {object} error
// @Router /deleteAdmin/{adminId} [delete]
func DeleteAdmin(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")

	adminCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	adminId := c.Param("adminId")

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(adminId) 
	result, err := adminCollection.DeleteOne(ctx, bson.M{"_Id": objId}) 
	res := map[string]interface{}{"data": result}

	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		return
	}
	
	if result.DeletedCount < 1 {
		logs.Error("No data to delete")
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"DeleteAdmin.500"))
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": localization.GetMessage(languageToken,"201"), "Data": res})

}