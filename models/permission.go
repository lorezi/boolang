package models

type Permission struct {
	UUID    string   `json:"permission_id" bson:"permission_id,omitempty"`
	Type    string   `json:"type" bson:"type,omitempty"`
	Actions []string `json:"actions" bson:"actions,omitempty"`
}
