package user

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	"time"
	logs "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary return book from user 
// @ID user-book-return
// @Accept json
// @Produce json
// @Param uId header string true "UserID"
// @Param payload body model.UserBook true "Payload for User Book Return API"
// @Success 201 {object} model.User
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /UserBookReturn/ [patch]
func UserBooksReturn(c *gin.Context) {

	appsettingCollection := database.GetCollection("AppSetting")
	bookCollection := database.GetCollection("Books")
	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	t := time.Now().Format(time.Kitchen)
	d := time.Now().Weekday().String()


	settings := new(model.Timings) 

	id, _ := primitive.ObjectIDFromHex("63c1398db715e13c41025cb8")
	
	err := appsettingCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&settings)

	if err != nil {

		logs.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
		return

	}

	for _, day := range settings.Timing {
		
		if day.Day == d {

			if day.IsOpen {

				ns, _ := time.Parse(time.Kitchen, t)
				srt, _ := time.Parse(time.Kitchen, day.StartTime)
				end, _ := time.Parse(time.Kitchen, day.CloseTime)

				// fmt.Printf("ns: %v\n", ns)
				// fmt.Printf("srt: %v\n", srt)
				// fmt.Printf("end: %v\n", end)

				if !ns.After(srt) || !ns.Before(end) {

					logs.Error (day.Day+"timings:"+ "from -"+day.StartTime+"to"+day.CloseTime)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{ day.Day+"timings": "from:"+day.StartTime+"-"+day.CloseTime})
					return

				}  

			} else {

				logs.Error("library is closed on:", day.Day)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"library is closed on:": day.Day})
				return

			}
		}

	}

	userbook := new(model.UserBook)
	
	defer cancel()
	
	if err:= c.BindJSON(&userbook); err != nil {
		logs.Error(err.Error())
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
			logs.Error(err3)
			c.JSON(http.StatusNotFound, gin.H{"message": err3})
			return
		}

		userId := uId.(string)
		objId, _ := primitive.ObjectIDFromHex(userId)

		fmt.Printf("objId: %v\n", objId)
	

		err1 := userCollection.FindOne(ctx, bson.M{"_Id": objId, "Login": true}).Decode(&result1)
		err2 := bookCollection.FindOne(ctx, bson.M{"Title": userbook.Title}).Decode(&result2)
		
		if err1 != nil {
			logs.Error(err1.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not find user_id in user or user not logged in"})
			return
		}

		if err2 != nil {
			logs.Error(err2.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not find title in books"})
			return
		}

		bookstaken := result1.BooksTaken

		if len(bookstaken) == 0 {
			logs.Error("user currently has no books")
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

			logs.Error("book not available with user")
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
	
			logs.Error(err.Error())
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
				
				logs.Error("cannot update usercollection")
			
			}

		}

		c.JSON(http.StatusCreated, gin.H{"message": "booksreturn updated successfully!", "Data": booksreturn})
		
	}
}