package request

import "gopkg.in/mgo.v2/bson"

type DeleteRequest struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}
