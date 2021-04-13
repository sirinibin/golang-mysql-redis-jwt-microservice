package models

type Response struct {
	Status    bool              `bson:"status" json:"status"`
	Criterias interface{}       `bson:"criterias,omitempty" json:"criterias,omitempty"`
	Result    interface{}       `bson:"result,omitempty" json:"result,omitempty"`
	Errors    map[string]string `bson:"errors,omitempty" json:"errors,omitempty"`
}
