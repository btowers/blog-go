package adder

type Post struct {
	Author User   `bson:"author" json:"author"`
	Image  string `bson:"image" json:"image"  binding:"required"`
	Text   string `bson:"text" json:"text" binding:"required"`
}
