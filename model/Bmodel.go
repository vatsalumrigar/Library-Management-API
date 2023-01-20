package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Books struct {

	ID primitive.ObjectID `json:"_id" bson:"_Id"`
	Title string `json:"title" bson:"Title" binding:"required"`
	Author []Authors `json:"author" bson:"Author" binding:"required"`
	Description string `json:"description" bson:"Description" binding:"required"`
	Publisher Publishers `json:"publisher" bson:"Publisher" binding:"required"`
	Genre string `json:"genre" bson:"Genre" binding:"required"`
	Quantities int `json:"quantities" bson:"Quantities"`
	Status string `json:"status" bson:"Status" binding:"required"`
    Penalty int `json:"penalty" bson:"Penalty"`
	Cost int `json:"cost" bson:"Cost"`
	
}

type Authors struct {

	Name string `json:"name" bson:"Name" binding:"required"`
	Education string `json:"education" bson:"Education" `
	Author_Email string `json:"author_email" bson:"Author_Email"`

}

type Publishers struct {

	Companyname string `json:"company_name" bson:"Company_Name" binding:"required"`
	Owner string `json:"owner" bson:"Owner" binding:"required"`
	Publisher_Email string `json:"publisher_email" bson:"Publisher_Email"`
	PublishedOn string `json:"published_on" bson:"Published_On" binding:"required"`

}

type FilterModel struct {

	Title string `form:"title"`
	Genre string `form:"genre"`
	Author string `form:"author"`
	Publisher string `form:"publisher"`
	
}

type BooksIssued struct {

	BookID string `json:"bookid" bson:"BookID"`
	BookTitle string `json:"booktitle" bson:"BookTitle"`
	Cost int `json:"cost" bson:"Cost"`
	IssuedTo []IssueDetails `json:"issuedto" bson:"IssuedTo"`
	IssuedQuantity int `json:"issuedquantity" bson:"IssuedQuantity"`
	BooksLeft int `json:"booksleft" bson:"BooksLeft"`

}

type IssueDetails struct {

	UserID string `json:"userid" bson:"UserID"`
	Email string `json:"email" bson:"Email"`
	Quantity int `json:"quantity" bson:"Quantity"`

}

type HistoryPayload struct {

	BookTitle string `json:"booktitle" bson:"BookTitle"`

}

type Bookqty struct {
	Books map[string]interface{}  `json:"books" bson:"Books" binding:"required"`
	Operations string `json:"operations" bson:"Operations" binding:"required"`
}
