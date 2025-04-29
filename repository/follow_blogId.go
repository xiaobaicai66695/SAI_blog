package repository

type FollowBlog struct {
	Fid    int64 `gorm:"fid" json:"fid"`
	BlogId int64 `gorm:"blogId" json:"blogId"`
}
