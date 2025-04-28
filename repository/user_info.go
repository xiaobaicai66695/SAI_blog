package repository

type UserInfo struct {
	ID       int64  `gorm:"column:id"`
	Name     string `json:"name"`
	IcoUrl   string `json:"ico_url"`
	Follower string `json:"follower"`
	Followed string `json:"followed"`
}

func (UserInfo) TableName() string {
	return "user"
}

func QueryById(id int64) *UserInfo {
	var user UserInfo
	db.First(&user, id)
	return &user
}
