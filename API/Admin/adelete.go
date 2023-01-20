package admin

import (
	database "PR_2/databases"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
)

// @Summary delete admin
// @ID delete-admin
// @Produce json
// @Param adminId path string true "AdminID" 
// @Success 201 {object} model.User
// @Failure 500 {object} error
// @Router /deleteAdmin/{adminId} [delete]
func DeleteAdmin(c *gin.Context) {

	adminCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	adminId := c.Param("adminId")

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(adminId) 
	result, err := adminCollection.DeleteOne(ctx, bson.M{"_Id": objId}) 
	res := map[string]interface{}{"data": result}

	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	
	if result.DeletedCount < 1 {
		logs.Error("No data to delete")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No data to delete"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Admin deleted successfully", "Data": res})

}