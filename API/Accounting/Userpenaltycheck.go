package accounting

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary check penalty of user
// @ID acc-penalty-check
// @Accept json
// @Produce json
// @Param librarianId header string true "LibrarianID"
// @Param payload body model.Payload true "Query Payload for Penalty Check API"
// @Success 201 {object} model.Accounting
// @Failure 400 {object} error
// @Failure 403 {object} error
// @Failure 404 {object} error
// @Failure 406 {object} error
// @Failure 500 {object} error
// @Router /Accounting/penaltycheck [post]
func AccountingPenaltyCheck(c *gin.Context) {

	accountingCollection := database.GetCollection("Accounting")
	userCollection := database.GetCollection("User")
	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)

	defer cancel()
	payload := new(model.Payload)

	if err:= c.BindJSON(&payload); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return

	}

	if middleware.Authentication(c){

		lId,err := c.Get("uid")

		if !err {
			c.JSON(http.StatusNotFound, gin.H{"message": err})
			return
		}

		libId := lId.(string)
		objId1, _ := primitive.ObjectIDFromHex(libId)

		var lib model.User

		err1 := userCollection.FindOne(ctx, bson.M{"_Id": objId1}).Decode(&lib)
		
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err1})
			return
		}

		if lib.UserType != "Librarian"{
			c.JSON(http.StatusForbidden, gin.H{"message": "enter valid librairian token"})
			return
		}

		var user model.User

		objId2, _ := primitive.ObjectIDFromHex(payload.UserId)

		err2 := userCollection.FindOne(ctx, bson.M{"_Id": objId2}).Decode(&user)
		
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot find user"})
			return
		}

		for _, book := range payload.PenaltyDetail {

			bookfound := false

			if (book.Reason != 1) && (book.Reason != 2) && (book.Reason != 3) {

				c.AbortWithStatusJSON(http.StatusNotAcceptable,gin.H{"error":"reason should be either 1, 2 or 3"})
				return

			}		

			for _, ubook := range user.BooksTaken{
				
				if ubook.Title == book.BookTitle {
					bookfound = true
				}

			}

			if !bookfound {
				c.AbortWithStatusJSON(http.StatusForbidden,gin.H{"user has no book called" : book.BookTitle})
				return
			}

		}

		count_user, _ := accountingCollection.CountDocuments(ctx, bson.M{"UserId": payload.UserId})
		
		if len(user.BooksTaken) == 0{
			c.AbortWithStatusJSON(http.StatusNotFound,gin.H{"message": "user has no books"})
			return	
		}

		payloadpenlatydetail := payload.PenaltyDetail
		
		var pdetails []model.Pdetails	

		for ind, ubook := range user.BooksTaken {

			ubooktitle := ubook.Title
			bookmatch := false
			penalty := 0

			for _, book := range payloadpenlatydetail{
				
				if ubooktitle == book.BookTitle {

					
					pbooktitle := book.BookTitle
					reason := book.Reason

					var books model.Books

					err3 := bookCollection.FindOne(ctx, bson.M{"Title": pbooktitle}).Decode(&books)
		
					if err3 != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot find book"})
						return
					}

					// if reason == 1 {

					// 	reasontype := "Late Penalty"

					// 	timenow := time.Now().Unix()
					// 	timebooktaken := ubook.TimePenaltyCalc

					// 	tn := time.Unix(timenow,0)
					// 	tbt := time.Unix(timebooktaken,0)

					// 	diff := tn.Sub(tbt)
					// 	daydiff := int(diff.Hours()/24)

					// 	if daydiff == 0 {

					// 		penalty = 0

					// 	}

					// 	penaltycheckdays := 15

					// 	if ubook.TimePenaltyCalc != ubook.TimeTaken {

					// 		penaltycheckdays = 0

					// 	}

					// 	if daydiff > penaltycheckdays {

					// 		penalty = books.Penalty * (daydiff-penaltycheckdays)

					// 		penaltydetails := model.Pdetails{
					// 			LibrarianId: objId1.Hex(),     
					// 			BookTitle: ubook.Title,      
					// 			TimePenaltyCheck: time.Now().Unix(), 
					// 			PenaltyPay: false ,       
					// 			PenaltyAmount: penalty,    
					// 			Reason: reasontype,          
					// 			ReasonType: 1 ,       
					// 		}

					// 		fmt.Printf("penalty.1: %v\n", penalty)
					// 		pdetails = append(pdetails, penaltydetails)

					// 	}
					// }
					if reason == 1 {

					goto R1

				 	} else if reason == 2 {

						bookmatch = true
						reasontype := "Book Missing"

						penalty = books.Cost

						penaltydetails := model.Pdetails{
							LibrarianId: objId1.Hex(),     
							BookTitle: ubook.Title,      
							TimePenaltyCheck: time.Now().Unix(), 
							PenaltyPay: false ,       
							PenaltyAmount: penalty,    
							Reason: reasontype,          
							ReasonType: 2 ,      
						}

						fmt.Printf("penalty.2: %v\n", penalty)
						pdetails = append(pdetails, penaltydetails)

					} else if reason == 3 {

						bookmatch = true
						reasontype := "Book Damaged"

						penalty = (books.Cost*50)/100

						penaltydetails := model.Pdetails{
							LibrarianId: objId1.Hex(),     
							BookTitle: ubook.Title,      
							TimePenaltyCheck: time.Now().Unix(), 
							PenaltyPay: false ,       
							PenaltyAmount: penalty,    
							Reason: reasontype,          
							ReasonType: 3 ,      
						}

						fmt.Printf("penalty.3: %v\n", penalty)
						pdetails = append(pdetails, penaltydetails)

					} 

					
					match := bson.M{"_Id" : objId2}
					fmt.Printf("objId2: %v\n", objId2)
					indtostr := strconv.Itoa(ind)
					update := bson.M{
							
						"$set": bson.M{
							"Status": "Unavailable",
							"Books_Taken."+ indtostr + ".TimePenaltyCalc" : time.Now().Unix() ,
						},
						
					}
	
					_, err := userCollection.UpdateOne(ctx,match,update)

					if err != nil {
						
						c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user status or timepenaltycalc"})
						return
	
					}

				} 
			}

			R1 : if !bookmatch {

				var absentbook model.Books

				err6 := bookCollection.FindOne(ctx, bson.M{"Title": ubooktitle}).Decode(&absentbook)
		
				if err6 != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot find absentbook"})
					return
				}

				reason := 1
				reasontype := "Late Penalty"

				timenow := time.Now().Unix()
				timebooktaken := ubook.TimePenaltyCalc

				tn := time.Unix(timenow,0)
				tbt := time.Unix(timebooktaken,0)

				diff := tn.Sub(tbt)
				daydiff := int(diff.Hours()/24)

				if daydiff == 0 {

					penalty = 0

				}

				penaltycheckdays := 15

				if ubook.TimePenaltyCalc != ubook.TimeTaken {

					penaltycheckdays = 0

				}

				if daydiff > penaltycheckdays {

					penalty = absentbook.Penalty * (daydiff-penaltycheckdays)

					penaltydetails := model.Pdetails{
						LibrarianId: objId1.Hex(),     
						BookTitle: ubook.Title,      
						TimePenaltyCheck: time.Now().Unix(), 
						PenaltyPay: false ,       
						PenaltyAmount: penalty,    	
						Reason: reasontype,          
						ReasonType: reason ,       
					}

					fmt.Printf("penalty.1: %v\n", penalty)
					pdetails = append(pdetails, penaltydetails)
					fmt.Printf("pdetails1: %v\n", pdetails)

				}

				match := bson.M{"_Id" : objId2}
				fmt.Printf("objId2: %v\n", objId2)
				indtostr := strconv.Itoa(ind)
				update := bson.M{
						
					"$set": bson.M{
						"Status": "Unavailable",
						"Books_Taken."+ indtostr + ".TimePenaltyCalc" : time.Now().Unix() ,
					},
					
				}

				_, err := userCollection.UpdateOne(ctx,match,update)

				if err != nil {
					
					c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user status or timepenaltycalc"})
					return

				}
			}
		}

		totalpenalty := 0

		for _ , bk := range pdetails {

			totalpenalty = totalpenalty + bk.PenaltyAmount

		}

		match := bson.M{"_Id" : objId2}
		update := bson.M{
				
			"$set": bson.M{
				"Total_Penalty": totalpenalty ,
			},
			
		}

		_, err4 := userCollection.UpdateOne(ctx,match,update)

		if err4 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user total penalty"})
			return

		}

		accounting := model.Accounting {

			UserId: user.ID.Hex(),
			Firstname: user.Firstname,
			LastName: user.Lastname,
			Email: user.Email,
			TotalPenalty: totalpenalty,
			PenaltyDetail: pdetails,
			//TimePenaltyPay: ,
	
		}

		if count_user == 0 {
			
		_, err5 := accountingCollection.InsertOne(ctx, accounting)
		//res := map[string]interface{}{"data": result}
	
		if err5 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message":err5 })
			return
		}

		}

		if count_user > 0 {

			fmt.Printf("pdetails: %v\n", pdetails)
			match1 := bson.M{"UserId" : payload.UserId}
			//update1 := bson.M{"$push": bson.M{"PenaltyDetails": pdetails}}
			update1 := bson.M{"$addToSet": bson.M{"PenaltyDetails":	bson.M{"$each": pdetails}}}

			_, err6 := accountingCollection.UpdateOne(ctx,match1,update1)

			if err6 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "could not push books in pdetails"})
				return

			}

		}

		c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully"})
		return

	}
}


