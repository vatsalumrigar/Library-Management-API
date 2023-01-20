package admin

import (
	//validation "PR_2/validation"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary update admin
// @ID update-admin
// @Accept json
// @Produce json
// @Success 201 {object} model.User
// @Param adminId path string true "AdminID" 
// @Param payload body model.Admin true "Query Payload for update Admin API"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /updateAdmin/{adminId} [put]
func UpdateAdmin(c *gin.Context) {
	
	adminCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	
	
	adminId := c.Param("adminId")
	var admin model.Admin
	
	defer cancel()
	
	objId, _ := primitive.ObjectIDFromHex(adminId)
	
	if err := c.BindJSON(&admin); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	edited := bson.M {

		"UserType" : admin.UserType,
		"Firstname" : admin.Firstname,
		"Lastname" : admin.Lastname,
		"Email" : admin.Email ,
		"MobileNo" : admin.MobileNo ,
		"Password" : admin.Password ,
		"Username" : admin.Username ,
		"Status" : admin.Status,
		"Dob" : admin.Dob,
		"Login": admin.Login,
		"Address": admin.Address,

	}
	
	/*val := validation.ValidateUmodel(ctx, admin.Email, admin.Username, admin.MobileNo, admin.Dob, admin.Status)

	if val != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": val.Error() })
		return
	} */

	result, err := adminCollection.UpdateOne(ctx, bson.M{"_Id": objId}, bson.M{"$set": edited})

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

