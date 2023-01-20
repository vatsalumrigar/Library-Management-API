package user

import (
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
)

// @Summary show books taken by user
// @ID user-param
// @Produce json
// @Param user_id path string true "UserID"
// @Success 200 {object} model.ParamUser
// @Failure 404 {object} error
// @Router /UserParam/ [get]
func UserParam(c *gin.Context) {


	user_id := c.Request.URL.Query().Get("user_id")
	fmt.Println("user_id : ", user_id)

	userCollection := database.GetCollection("User")
	bookCollection := database.GetCollection("Books")
	ctx,cancel := database.DbContext(10)

	defer cancel()

	match := bson.M{}

	if user_id != "" {

		objid, _ := primitive.ObjectIDFromHex(user_id) 
		match = bson.M{"_Id": objid}

	}

	var userparam []model.ParamUser

	cursor, err := userCollection.Find(ctx, match)
	
	//count_title, _ := userCollection.CountDocuments(ctx, match)

	//fmt.Printf("count_title: %v\n", count_title)

	if err != nil {

		logs.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err})
		return
		
	}

	for cursor.Next(ctx){	

		var founduser model.User

		err := cursor.Decode(&founduser)

		if err!= nil {
			logs.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err})
			return
		}
	
		var books model.Books
		var bookstaken []model.Bookdetail2

		for _, book := range founduser.BooksTaken {
			
			err := bookCollection.FindOne(ctx, bson.M{"Title": book.Title}).Decode(&books)
			
			if err!= nil {
				logs.Error(err.Error())
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err})
				return
			} 

			var qtyupdate bool
			qtyupdate = false

			qty := 1
			updateIndex := 0

			for index, booktitle := range bookstaken  {
				
				if book.Title == booktitle.Title {
	
					updateIndex = index
				
					qtyupdate = true

				}

			}	

			if !qtyupdate {

				btaken := model.Bookdetail2{
					Title : books.Title,
					Author : books.Author,
					Publisher : books.Publisher,
					Quantities : qty,
				} 
	
				bookstaken = append(bookstaken, btaken)

			} else {
				bookstaken[updateIndex].Quantities += 1
			}

		}
		
		uparam := model.ParamUser{
			User_Id : founduser.ID.Hex(),
			Username : founduser.Username,
			Email : founduser.Email,
			BookTaken : bookstaken,
		}

		userparam = append(userparam, uparam)

	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Posted successfully", "Data": userparam})	

}


