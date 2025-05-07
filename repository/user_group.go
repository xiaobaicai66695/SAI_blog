package repository

type UserGroup struct {
	Uid     int64 `gorm:"primary_key;AUTO_INCREMENT";json:"uid"`
	GroupId int64 `json:"group_id"`
}
type Group struct {
	GroupId   int64  `gorm:"primary_key;AUTO_INCREMENT"`
	FromUid   int64  `json:"from_uid"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
