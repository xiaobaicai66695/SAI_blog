package repository

type BlogComment struct {
	BlogId  int64  `gorm:"column:blog_id" json:"blog_id"`
	Cid     int64  `gorm:"column:cid" json:"cid"`
	Comment string `gorm:"column:comment;type:text" json:"comment"`
}

func QueryCommentsById(blogId int64) []BlogComment {
	var comments []BlogComment
	db.Where("blog_id = ?", blogId).Find(&comments)
	return comments
}
