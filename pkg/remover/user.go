package remover

type User struct {
	Email string `bson:"email"  binding:"required"`
}
