package repository

import "gorm.io/gorm"

type FollowBlog struct {
	Uid    int64 `gorm:"uid" json:"uid"`
	BlogId int64 `gorm:"blogId" json:"blogId"`
}

func (FollowBlog) TableName() string {
	return "follow_blogs"
}

func AddFollowBlog(tx *gorm.DB, uid int64, fid int64) {
	var blogs []Blog
	//里面的uid是作者的id，
	tx.Where("uid = ?", fid).Find(&blogs)
	for _, b := range blogs {
		var followBlog = &FollowBlog{
			Uid:    uid,
			BlogId: b.BlogId,
		}
		db.Create(followBlog)
	}
}

func DelFollowBlog(tx *gorm.DB, fid int64) {
	var blogs []Blog
	//里面的uid是作者的id，
	db.Where("uid = ?", fid).Find(&blogs)
	for _, b := range blogs {
		tx.Where("blog_id = ?", b.BlogId).Delete(&FollowBlog{})
	}
}
