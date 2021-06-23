package auth

type User struct {
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Email     string `bson:"email" json:"email"  binding:"required"`
	Password  string `bson:"password" json:"password" binding:"required"`
}
