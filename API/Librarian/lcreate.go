package librarians

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	logs "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary create librarian
// @ID create-librarian
// @Accept json
// @Produce json
// @Param payload body model.User true "Payload for create Librarian API"
// @Success 201 {object} model.User
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /Librarian/ [post]
func CreateLibrarian(c *gin.Context){

	librarianCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	librarian := new(model.User)

	defer cancel()

	if err := c.BindJSON(&librarian); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	addedLibrarian := model.User {
			
		ID: primitive.NewObjectID(),
		UserType: librarian.UserType,
		Firstname : librarian.Firstname,
		Lastname : librarian.Lastname,
		Email : librarian.Email,
		MobileNo : librarian.MobileNo,
		Password : librarian.Password,
		Username : librarian.Username,
		Status : librarian.Status,
		Dob : librarian.Dob,
		Login : librarian.Login,
		Address : librarian.Address,

	}

	result, err := librarianCollection.InsertOne(ctx, addedLibrarian)
	res := map[string]interface{}{"data": result}

	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error() })
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": res})

}