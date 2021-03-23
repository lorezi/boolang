package models

type PermissionGroup struct {
	GroupName  string       `json:"group_name" bson:"group_name"`
	GroupID    string       `json:"group_id" bson:"group_id"`
	Permission []Permission `json:"permission" bson:"permission"`
}

type Permission struct {
	Role    string    `json:"role" bson:"role"`
	Actions []Actions `json:"actions" bson:"actions"`
}

type Actions struct {
	Create bool `json:"create" bson:"create"`
	Read   bool `json:"read" bson:"read"`
	Update bool `json:"update" bson:"update"`
	Delete bool `json:"delete" bson:"delete"`
}
