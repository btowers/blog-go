package lister

type Post struct {
	Id     string `bson:"_id" json:"_id" binding:"required"`
	Author User   `bson:"author" json:"author"`
	Image  string `bson:"image" json:"image"  binding:"required"`
	Text   string `bson:"text" json:"text" binding:"required"`
}
