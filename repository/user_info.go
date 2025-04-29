package repository

type UserInfo struct {
	ID       int64  `gorm:"column:id"`
	Name     string `json:"name"`
	IcoUrl   string `json:"ico_url"`
	Follower string `json:"follower"`
	Followed string `json:"followed"`
	IsFollow bool   `json:"is_follow"`
}

func (UserInfo) TableName() string {
	return "user"
}

func QueryById(id int64) *UserInfo {
	var user UserInfo
	db.First(&user, id)
	err := db.First(Relaationship{Uid: id}).Error
	if err != nil {
		user.IsFollow = false
	}
	user.IsFollow = true
	return &user
}
