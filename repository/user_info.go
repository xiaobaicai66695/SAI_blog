package repository

type UserInfo struct {
	ID       int64  `gorm:"column:id"`
	Name     string `json:"name"`
	IcoUrl   string `json:"ico_url"`
	Follower int64  `json:"follower"`
	Followed int64  `json:"followed"`
	IsFollow bool   `json:"is_follow"`
}

func (UserInfo) TableName() string {
	return "user"
}

func QueryUserInfoById(id int64, uid int64) *UserInfo {
	var user UserInfo
	db.First(&user, id)
	var rel Relationship
	db.Where("uid = ? and fid = ?", uid, id).First(&rel)
	if uid == 0 || rel.Fid == id {
		user.IsFollow = true
	} else {
		user.IsFollow = false
	}

	return &user
}
