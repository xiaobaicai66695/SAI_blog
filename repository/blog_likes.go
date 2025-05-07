package repository

type BlogLikes struct {
	BlogId int64 `gorm:"column:blog_id" json:"blog_id"`
	Lid    int64 `gorm:"column:lid" json:"lid"`
}

func QueryLikesById(blogId int64) []BlogLikes {
	var likes []BlogLikes
	db.Where("blog_id=?", blogId).Find(&likes)
	return likes
}
