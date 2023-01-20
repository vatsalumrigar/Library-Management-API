package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)



type User struct {

	ID primitive.ObjectID `json:"_id" bson:"_Id" `
	UserType string `json:"user_type" bson:"User_Type"`
	Firstname string `json:"first_name" bson:"Firstname"`
	Lastname string `json:"last_name" bson:"Lastname"`
	Email string `json:"email" bson:"Email" binding:"required"`
	MobileNo string `json:"mobile_no" bson:"Mobile_No"`
	Password string `json:"password" bson:"Password"`
	Username string `json:"username" bson:"Username"`
	BooksTaken []Bookdetails `json:"books_taken" bson:"Books_Taken"`
	Status string `json:"status" bson:"Status"`
	Dob string `json:"dob" bson:"Dob"`
	Login bool `json:"login" bson:"Login"`
	IsFirstLogin bool `json:"isfirstlogin" bson:"IsFirstLogin"`
	Total_Penalty int `json:"total_penalty" bson:"Total_Penalty"`
	Address Address `json:"address" bson:"Address"`
	
}


type Token struct {

	Token string `json:"token" bson:"Token"`
    Refresh_Token string `json:"refresh_token" bson:"Refresh_Token"`
    Created_at time.Time `json:"created_at"  bson:"Created_At"`
    Updated_at time.Time `json:"updated_at"  bson:"Updated_At"`
    User_id string `json:"user_id"  bson:"User_Id"`

}

type Bookdetails struct {

	BookId string `json:"book_id" bson:"Book_Id"`
	Title string `json:"title" bson:"Title" binding:"required"`
	TimeTaken int64 `json:"time_taken" bson:"Time_Taken" `
	TimePenaltyCalc int64 `json:"timepenaltycalc" bson:"TimePenaltyCalc"`

}


type UserBook struct {

	// User_Id string `json:"user_id" bson:"User_Id" binding:"required"`
	Title string	`json:"title" bson:"Title" binding:"required"`

 }

type Login struct {

	Email string `json:"email" bson:"Email" binding:"required"`	
	Password string `json:"password" bson:"Password" binding:"required"`
	
}

type PenaltyUsers struct {

	User_id []string `json:"user_id" bson:"User_Id" binding:"required"`

}

type Logout struct {

	Email string `json:"email" bson:"Email" binding:"required"`
	
}

type IsPenalty struct {

	Username string `json:"username"`
	Bookname []string `json:"bookname"`
	Penalty int `json:"penalty"`

}

type PenaltyPay struct {

	Username string `json:"username" bson:"Username" binding:"required"`
	Pay_Amount int `json:"pay_amount" bson:"Pay_Amount" binding:"required"`
	
}

type ParamPayload struct {

	User_Id string `json:"user_id" bson:"User_Id" binding:"required"`

}

type ParamUser struct {

	User_Id string `json:"user_id" bson:"User_Id" binding:"required"`
	Username string `json:"username" bson:"Username" binding:"required"`
	Email string `json:"email" bson:"Email" binding:"required"`
	BookTaken []Bookdetail2 `json:"books_taken" bson:"Books_Taken"`

}

type Bookdetail2 struct {

	Title string `json:"title" bson:"Title" binding:"required"`
	Author []Authors `json:"author" bson:"Author"`
	Publisher Publishers `json:"publisher" bson:"Publisher" binding:"required"`
	Quantities int `json:"quantities" bson:"Quantities"`
	
}

type SetNewPassword struct {

	Email string `json:"email" bson:"Email" binding:"required"`	
	OldPassword string `json:"oldpassword" bson:"OldPassword" binding:"required"`
	NewPassword string `json:"newpassword" bson:"NewPassword" binding:"required"`

}
