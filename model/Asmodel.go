package model

type Timings struct {
	Timing []Times `json:"timing" bson:"Timing"`
}

type Times struct {
	Day       string `json:"day" bson:"Day"`
	IsOpen    bool   `json:"isopen" bson:"IsOpen"`
	StartTime string `json:"starttime" bson:"StartTime"`
	CloseTime string `json:"closetime" bson:"CloseTime"`
}