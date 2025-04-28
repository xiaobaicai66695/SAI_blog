package repository

type BlogLikes struct {
	BlogId int64 `gorm:"column:blog_id" json:"blog_id"`
	Lid    int64 `gorm:"column:lid" json:"lid"`
}
