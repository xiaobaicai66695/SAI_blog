package repository

type FollowBlog struct {
	Uid    int64 `gorm:"uid" json:"uid"`
	BlogId int64 `gorm:"blogId" json:"blogId"`
}

func (FollowBlog) TableName() string {
	return "follow_blogs"
}

func AddFollowBlog(uid int64, fid int64) {
	var blogs []Blog
	//里面的uid是作者的id，
	db.Where("uid = ?", fid).Find(&blogs)
	for _, b := range blogs {
		var followBlog = &FollowBlog{
			Uid:    uid,
			BlogId: b.BlogId,
		}
		db.Create(followBlog)
	}
}

func DelFollowBlog(uid int64, fid int64) {
	var blogs []Blog
	//里面的uid是作者的id，
	db.Where("uid = ?", fid).Find(&blogs)
	for _, b := range blogs {
		db.Where("blog_id = ?", b.BlogId).Delete(&FollowBlog{})
	}
}
