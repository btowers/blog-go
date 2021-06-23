package updater

type User struct {
	FirstName string `bson:"firstName" json:"firstName"  binding:"required"`
	LastName  string `bson:"lastName" json:"lastName"  binding:"required"`
	Email     string `bson:"email" json:"email"  binding:"required"`
	Password  string `bson:"password" json:"password" binding:"required"`
}
