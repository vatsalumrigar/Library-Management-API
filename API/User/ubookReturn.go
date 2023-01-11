package user

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func UserBooksReturn(c *gin.Context) {

	bookCollection := database.GetCollection("Books")
	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)


	userbook := new(model.UserBook)
	
	defer cancel()
	
	if err:= c.BindJSON(&userbook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	
	if middleware.Authentication(c) {

		var result1 model.User
		var result2 model.Books

		// userId := c.Param("userId")
		// objId, _ := primitive.ObjectIDFromHex(userId)

		uId,err3 := c.Get("uid")

		if !err3 {
			c.JSON(http.StatusNotFound, gin.H{"message": err3})
			return
		}

		userId := uId.(string)
		objId, _ := primitive.ObjectIDFromHex(userId)

		fmt.Printf("objId: %v\n", objId)
	

		err1 := userCollection.FindOne(ctx, bson.M{"_Id": objId, "Login": true}).Decode(&result1)
		err2 := bookCollection.FindOne(ctx, bson.M{"Title": userbook.Title}).Decode(&result2)
		
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not find user_id in user or user not logged in"})
			return
		}

		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not find title in books"})
			return
		}

		bookstaken := result1.BooksTaken

		if len(bookstaken) == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{"error":"user currently has no books"})
			return
		}

		var notfound bool
		notfound = true
		qty1 := 0

		for _ , title := range bookstaken {
			
			giventitle := userbook.Title		
			tomatchtitle := title.Title

			if giventitle == tomatchtitle{
				qty1 += 1
				notfound = false
			}

		}

		if notfound {

			c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{"error":"book not available with user"})
			return
			
		}

		match := bson.M{"_Id": objId}

		var bd = model.Bookdetails{

					BookId : result2.ID.Hex(),
					Title: result2.Title,
					TimeTaken: time.Now().Unix(),

		}

		change := bson.M{"$pull": bson.M{"Books_Taken": bson.M{"Title":bd.Title}}}
		fmt.Printf("change: %v\n", change)


		booksreturn, err := userCollection.UpdateOne(ctx, match, change)
	
		if err != nil {
	
			c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{"error":"cannot update usercollection"})
			return

		} else {


			qty := result2.Quantities + qty1

			update := bson.M {
				"$set": bson.M {
					"Quantities" : qty,
				},
			}

			if qty > 0 {
				update = bson.M {
					"$set": bson.M {
						"Quantities" : qty,
						"Status" : "Available",
					},
				}
			}

			_ , err := bookCollection.UpdateOne(ctx, bson.M{"Title": result2.Title}, update)

			if err!= nil {

				fmt.Println("cannot update usercollection")
			}

		}

		c.JSON(http.StatusCreated, gin.H{"message": "booksreturn updated successfully!", "Data": booksreturn})
		
	}
}