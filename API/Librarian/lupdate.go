package librarians

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary update librarian
// @ID update-librarian
// @Accept json
// @Produce json
// @Param librarianId path string true "LibrarianID" 
// @Param payload body model.User true "Payload for update Librarian API"
// @Success 201 {object} model.User
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /updateLibrarian/{librarianId} [put]
func UpdateLibrarian(c *gin.Context) {
	
	librarianCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	
	
	librarianId := c.Param("librarianId")
	var librarian model.User
	
	defer cancel()
	
	objId, _ := primitive.ObjectIDFromHex(librarianId)
	
	if err := c.BindJSON(&librarian); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	edited := bson.M {

		"UserType" : librarian.UserType,
		"Firstname" : librarian.Firstname,
		"Lastname" : librarian.Lastname,
		"Email" : librarian.Email ,
		"MobileNo" : librarian.MobileNo ,
		"Password" : librarian.Password ,
		"Username" : librarian.Username ,
		"Status" : librarian.Status,
		"Dob" : librarian.Dob,
		"Login": librarian.Login,
		"Address": librarian.Address,

	}
	
	/*val := validation.ValidateUmodel(ctx, admin.Email, admin.Username, admin.MobileNo, admin.Dob, admin.Status)

	if val != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": val.Error() })
		return
	} */

	result, err := librarianCollection.UpdateOne(ctx, bson.M{"_Id": objId}, bson.M{"$set": edited})

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

