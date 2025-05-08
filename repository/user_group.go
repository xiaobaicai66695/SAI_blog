package repository

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type UserGroup struct {
	Uid int64 `json:"uid"`
	//0 web组 1 AI组 2 安卓组 3 游戏组
	GroupId int64 `json:"group_id"`
}
type Group struct {
	GroupId   int64  `gorm:"primary_key;AUTO_INCREMENT"`
	FromUid   int64  `json:"from_uid"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

func QueryGroupIdByUid(uid int64, groupId int64) error {
	var group Group
	db.Where("uid=? and group_id", uid, groupId).First(&group)
	if group.GroupId == 0 {
		return fmt.Errorf("没有进入这个群聊")
	}
	return nil
}

func QueryUidByGroupId(groupId int64) []int64 {
	var uids []int64
	db.Where("group_id=?", groupId).Find(&uids)
	return uids
}

var groupUsers = make(map[int64][]*websocket.Conn)

func AddUserToGroup(userId int64, groupId int64, conn *websocket.Conn) error {
	//if groupUsers[userId] == nil {
	//	groupUsers[userId] = make([]*websocket.Conn, 0)
	//}
	//groupUsers[userId] = append(groupUsers[userId], conn)
	return nil
}

func GetUserIdsByGroupId(groupId int64) ([]int64, error) {
	return nil, nil
}

func SaveMsgToRedis(msg *Group, groupId int64) error {
	return nil
}
