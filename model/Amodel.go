package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {

	ID primitive.ObjectID `json:"_id" bson:"_Id" `
	UserType string `json:"user_type" bson:"User_Type" binding:"required"`
	Firstname string `json:"first_name" bson:"Firstname" binding:"required"`
	Lastname string `json:"last_name" bson:"Lastname" binding:"required"`
	Email string `json:"email" bson:"Email" binding:"required"`
	MobileNo string `json:"mobile_no" bson:"Mobile_No" binding:"required"`
	Password string `json:"password" bson:"Password" binding:"required"`
	Username string `json:"username" bson:"Username" binding:"required"`
	Status string `json:"status" bson:"Status"`
	Dob string `json:"dob" bson:"Dob" binding:"required"`
	Login bool `json:"login" bson:"Login" binding:"required"`
	Address Address `json:"address" bson:"Address"`

}


type Address struct {

	Street string `json:"street" bson:"Street"`
	City string `json:"city" bson:"City"`
	State string `json:"state" bson:"State"`
	Pincode int `json:"pincode" bson:"Pincode"`
	Country string `json:"country" bson:"Country"`

}