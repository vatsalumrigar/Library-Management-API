package book

import (
	validation "PR_2/validation"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
)

// @Summary update book
// @ID update-book
// @Accept json
// @Produce json
// @Param bookId path string true "BookID"
// @Param payload body model.Books true "Payload for update Books API"
// @Success 201 {object} model.Books
// @Failure 409 {object} error
// @Failure 500 {object} error
// @Router /updateBook/{bookId} [put]
func UpdateBook(c *gin.Context) {

	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)

	
	bookId := c.Param("bookId")
	var book model.Books
	
	

	defer cancel()
	
	objId, _ := primitive.ObjectIDFromHex(bookId)
	
	if err := c.BindJSON(&book); err != nil {

		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return

	}
	
	for _ , author := range book.Author {
		//fmt.Printf("%d : %s, ", i,  author.Author_Email)
		val := validation.ValidateBmodel(ctx, author.Author_Email, book.Publisher.Publisher_Email, book.Publisher.PublishedOn )
	
		if val != nil {
	
				logs.Error(val.Error())
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": val.Error() })
				return
		
			}
		}

	edited := bson.M {

		"Title" : book.Title,
		"Author" : book.Author,
		"Description" : book.Description,
		"Publisher" : book.Publisher,
		"Genre" : book.Genre,
		"Quantities" : book.Quantities,
		"Status" : book.Status,	
		"Penalty": book.Penalty,
		
	}
	
	

	//edited := bson.M{"title": post.Title, "article": post.Article}


	result, err := bookCollection.UpdateOne(ctx, bson.M{"_Id": objId}, bson.M{"$set": edited})

	res := map[string]interface{}{"data": result}

	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
		}

	if result.MatchedCount < 1 {
		logs.Error("Data doesn't exist")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Data doesn't exist"})
		return
	}
			
	c.JSON(http.StatusCreated, gin.H{"message": "data updated successfully!", "Data": res})	
}