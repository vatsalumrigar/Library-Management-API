package book

import (
	database "PR_2/databases"
	model "PR_2/model"
	validation "PR_2/validation"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateBook(c *gin.Context) {

	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)

	book := new(model.Books)
	//author := new(bmodel.Authors)
	
	defer cancel()

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	} 
	
	for _ , author := range book.Author {
	//fmt.Printf("%d : %s, ", i,  author.Author_Email)
	err := validation.ValidateBmodel(ctx, author.Author_Email, book.Publisher.Publisher_Email, book.Publisher.PublishedOn )

	if err != nil {

			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error() })
			return
	
		}
	}
	//err := validateBmodel(ctx, book.Author[0].Email, book.Publisher.Email )


	//appendBookQty(ctx, book.Title, book.Quantities)
	//bookId := c.Param("bookId")
	var result model.Books
	//objId, _ := primitive.ObjectIDFromHex(bookId)
	
	title := book.Title
	quantity := book.Quantities
	count_title, _ := bookCollection.CountDocuments(ctx , bson.M{"Title": title} )

	if count_title >= 1 {

		get1 := bookCollection.FindOne(ctx, bson.M{"Title": title}).Decode(&result)
		if get1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": get1})
			return
		}
	
		qty := quantity + result.Quantities
		update := bson.M{
			
				"$set": bson.M{
					"quantities": qty ,
				},
			
		}
		
		c.JSON(http.StatusCreated, gin.H{"message": "qty updated successfully!", "Data": update})
		
		edited := bson.M {

			"Title" : book.Title,
			"Author" : book.Author,
			"Description" : book.Description,
			"Publisher" : book.Publisher,
			"Genre" : book.Genre,
			"Quantities" : qty,
			"Status" : book.Status,	
			"Penalty" : book.Penalty,
			
		}
		err, _ := bookCollection.UpdateOne(ctx,bson.M{"Title": title},bson.M{"$set": edited})
	
		if err != nil {
		
			return
		}

	} else {

		addedBook := model.Books{

			ID: primitive.NewObjectID(),
			Title: book.Title,
			Author: book.Author, 
			Description: book.Description,
			Publisher: book.Publisher,
			Genre: book.Genre,
			Quantities: book.Quantities,
			Status: book.Status,
			Penalty: book.Penalty,

		}

		result1, err := bookCollection.InsertOne(ctx, addedBook)

		res1 := map[string]interface{}{"data": result1}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": res1})

	}

}

