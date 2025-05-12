package repository

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

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

func UploadComment(blogId int64, uid int64, msg string) {
	key := fmt.Sprintf("blog:%d", blogId)
	db.Create(&BlogComment{
		BlogId:  blogId,
		Cid:     uid,
		Comment: msg,
	})
	blogJson, _ := rdb1.Get(ctx, key).Result()
	var blog Blog
	_ = json.Unmarshal([]byte(blogJson), &blog)
	blog.Comments++
	newBlogJson, _ := json.Marshal(blog)
	rdb1.Set(ctx, key, newBlogJson, time.Hour*10)
	db.Model(&Blog{}).Where("blog_id = ?", blogId).Update("comments", gorm.Expr("comments+?", 1))
	return
}
