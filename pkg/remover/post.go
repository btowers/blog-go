package remover

type Post struct {
	Id string `bson:"_id" json:"_id" binding:"required"`
}
