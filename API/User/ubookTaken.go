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

// @Summary add books to user 
// @ID user-book-take
// @Accept json
// @Produce json
// @Param uId header string true "UserID"
// @Param payload body model.UserBook true "Payload for User Book Taken API"
// @Success 201 {object} model.User
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /UserBookTaken/ [patch]
func UserBooksTaken(c *gin.Context) {

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
					logs.Error(day.Day+"timings:"+"from -"+day.StartTime+"to"+day.CloseTime)
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

	//book := new(bmodel.Books)
	//user := new(umodel.User)
	userbook := new(model.UserBook)
	
	defer cancel()

	if err:= c.BindJSON(&userbook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		logs.Error(err.Error())
		return
	}
	
	if middleware.Authentication(c) {
	
		// //if user "Available" -> count documents
		uId,err1 := c.Get("uid")
		fmt.Printf("uId: %v\n", uId)

		if !err1 {
			logs.Error(err1)
			c.JSON(http.StatusNotFound, gin.H{"message": err1})
			return
		}

		userId := uId.(string)
		//userId := c.Param("userId")
		fmt.Printf("userId: %v\n", userId)

		var res model.User	
		
		objId, _ := primitive.ObjectIDFromHex(userId)

		fmt.Printf("objId: %v\n", objId)


		err := userCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&res)


		if err != nil {
			logs.Error(err1)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "user not logged in"})
			return
		}

		if res.Status == "Available" {

			count_title, _ := bookCollection.CountDocuments(ctx, bson.M{"Title": userbook.Title})
			

			if count_title == 0 {

					logs.Error("could not find title in books")				
					c.JSON(http.StatusNotFound, gin.H{"message": "could not find title in books"})
					return
				
			} 

			var result model.Books

			if count_title >= 1 {

				err := bookCollection.FindOne(ctx, bson.M{"Title": userbook.Title}).Decode(&result)

				if err != nil {
					logs.Error(err.Error())
					c.JSON(http.StatusNotFound, gin.H{"message": "could not find title in books"})
					return
				}

				if result.Status == "Available" && result.Quantities > 0 {
					
				qty := result.Quantities - 1


				update := bson.M{
					
						"$set": bson.M{
							"Quantities": qty ,
						},
					
				}


				if qty == 0 {

					update = bson.M{
						"$set": bson.M {
							"Quantities": qty ,
							"Status": "Unavailable",
						},
					}
				}


				_, err := bookCollection.UpdateOne(ctx, bson.M{"Title": result.Title},update)

				if err != nil {
					logs.Error(err.Error())
					c.JSON(http.StatusNotFound, gin.H{"message": "could not update book" })
					return
				}

				// userId := c.Param("userId")
				// objId, _ := primitive.ObjectIDFromHex(userId)

				// uId,err1 := c.Get("uid")

				// if !err1 {
				// 	c.JSON(http.StatusNotFound, gin.H{"message": err1})
				// 	return
				// }

				// userId := uId.(string)
				// objId, _ := primitive.ObjectIDFromHex(userId)

			
				match := bson.M{"_Id" : objId}
				change := bson.M{"$push": bson.M{"Books_Taken": model.Bookdetails {
					BookId : result.ID.Hex(),
					Title: result.Title,
					TimeTaken: time.Now().Unix(),
					TimePenaltyCalc: time.Now().Unix(),
					
				} }}

				newbook, err:= userCollection.UpdateOne(ctx, match, change)

				if err !=  nil {
					logs.Error(err.Error())
					c.JSON(http.StatusNotFound, gin.H{"message": "could not update bookstaken" })
					return
				}

				c.JSON(http.StatusCreated, gin.H{"message": "bookstaken updated successfully!", "Data": newbook})	
				return
			
				} else {
					logs.Error("not enough book qty")
					c.JSON(http.StatusNotFound, gin.H{"message": "not enough book qty" })
					return
				}
				
			}

		} else {

			c.JSON(http.StatusNotFound, gin.H{"message": "user not available" })
			return

		}
	}else{

		logs.Error("user not logged in")
		c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{"message": "user not logged in" })
		return

	}

}	