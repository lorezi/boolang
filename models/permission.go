package models

type Permission struct {
	Type    string   `json:"type" bson:"type,omitempty"`
	Actions []string `json:"actions" bson:"actions,omitempty"`
}
