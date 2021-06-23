package lister

type User struct {
	Id        string `bson:"_id" json:"_id"  binding:"required"`
	FirstName string `bson:"firstName" json:"firstName"  binding:"required"`
	LastName  string `bson:"lastName" json:"lastName"  binding:"required"`
	Email     string `bson:"email" json:"email"  binding:"required"`
}
