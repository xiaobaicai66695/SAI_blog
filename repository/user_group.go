package repository

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"time"
)

type UserGroup struct {
	Uid int64 `json:"uid"`
	//0 web组 1 AI组 2 安卓组 3 游戏组
	GroupId int64 `json:"group_id"`
}
type Group struct {
	GroupId   int64  `json:"group_id"`
	FromUser  UserVO `json:"from_user"`
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
	var users []UserGroup
	err := db.Where("group_id=?", groupId).Find(&users).Error
	if err != nil {
		return nil, err
	}
	var uids []int64
	for _, user := range users {
		uid := user.Uid
		uids = append(uids, uid)
	}
	return uids, nil
}

func SaveMsgToRedis(msg *Group, groupId int64) error {
	key := fmt.Sprintf("msg:history:%d", groupId)
	msgStr, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("存入redis有问题，msg未能转化未json字符串")
	}
	rdb3.ZAdd(ctx, key, &redis.Z{
		Score:  float64(msg.Timestamp),
		Member: string(msgStr),
	})
	rdb3.Expire(ctx, key, time.Hour*1440)
	return nil
}
func QueryUidInGroupIsOrNot(uid int64, groupId int64) bool {
	var userGroup UserGroup
	err := db.Where("uid=? and group_id=?", uid, groupId).First(userGroup).Error
	if userGroup.GroupId == 0 || err != nil {
		return false
	}
	return true
}

func PullHistoryMsg(groupId int64, cursor int64) (msgs []*Group, err error) {
	key := fmt.Sprintf("msg:history:%d", groupId)
	msgStrs, err := rdb3.ZRevRange(ctx, key, cursor*50, (cursor+1)*50).Result()
	if err != nil {
		return nil, err
	}
	for _, msgStr := range msgStrs {
		var msg Group
		_ = json.Unmarshal([]byte(msgStr), &msg)
		msgs = append(msgs, &msg)
	}
	return msgs, nil
}
