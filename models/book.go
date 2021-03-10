package models

// ExpBook model
type ExpBook struct {
	ID   string `json:"id" bson:"_id"`
	Book `bson:"inline"`
}

// Book model for creating new book
type Book struct {
	Title  string `json:"title" bson:"title"`
	Author string `json:"author" bson:"author"`
	Year   string `json:"year" bson:"year"`
}
