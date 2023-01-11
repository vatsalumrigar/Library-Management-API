package model

type Accounting struct {
	UserId         string     `json:"userid" bson:"UserId"`
	Firstname      string     `json:"firstname" bson:"FirstName"`
	LastName       string     `json:"lastname" bson:"LastName"`
	Email          string     `json:"email" bson:"Email"`
	TotalPenalty   int        `json:"totalpenalty" bson:"TotalPenalty"`
	PenaltyDetail  []Pdetails `json:"penaltydetails" bson:"PenaltyDetails"`
	TimePenaltyPay int64      `json:"timepenaltypay" bson:"TimePenaltyPay"`
}

type Accounting2 struct {
	UserId         string     `json:"userid" bson:"UserId"`
	Firstname      string     `json:"firstname" bson:"FirstName"`
	LastName       string     `json:"lastname" bson:"LastName"`
	Email          string     `json:"email" bson:"Email"`
	TotalPenalty   int        `json:"totalpenalty" bson:"TotalPenalty"`
	PenaltyDetail  Pdetails `json:"penaltydetails" bson:"PenaltyDetails"`
	TimePenaltyPay int64      `json:"timepenaltypay" bson:"TimePenaltyPay"`
}

type Pdetails struct {
	LibrarianId      string
	BookTitle        string `json:"booktitle" bson:"BookTitle"`
	TimePenaltyCheck int64  `json:"timepenaltycheck" bson:"TimePenaltyCheck"`
	PenaltyPay       bool   `json:"penaltypay" bson:"PenaltyPay"`
	PenaltyAmount    int    `json:"penaltyamount" bson:"PenaltyAmount"`
	Reason           string `json:"reason" bson:"Reason"`
	ReasonType       int    `json:"reasontype" bson:"ReasonType"`
}

type Payload struct {
	UserId        string           `json:"userid" bson:"UserId" binding:"required"`
	PenaltyDetail []PenaltyDetails `json:"penaltydetail" bson:"PenaltyDetail" binding:"required"`
}

type PenaltyDetails struct {
	BookTitle string `json:"booktitle" bson:"BookTitle"`
	Reason    int    `json:"reason" bson:"Reason"`
}