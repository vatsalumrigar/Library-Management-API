// @title Library Management API
// @version 1.0
// @description This is a  Library Management API server.
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /
// @query.collection.format multi
package book

import (
	validation "PR_2/validation"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary update book
// @ID update-book
// @Accept json
// @Produce json
// @Param bookId path string true "BookID"
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	
	for _ , author := range book.Author {
		//fmt.Printf("%d : %s, ", i,  author.Author_Email)
		val := validation.ValidateBmodel(ctx, author.Author_Email, book.Publisher.Publisher_Email, book.Publisher.PublishedOn )
	
		if val != nil {
	
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
		}

	if result.MatchedCount < 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Data doesn't exist"})
		return
	}
			
	c.JSON(http.StatusCreated, gin.H{"message": "data updated successfully!", "Data": res})	
}