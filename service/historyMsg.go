package service

import (
	"SAI_blog/repository"
)

func PullHistoryMsg(groupId int64, cursor int64) ([]*repository.Group, error) {
	//ok := repository.QueryUidInGroupIsOrNot
	//if !ok {
	//	fmt.Errorf("不在群聊内")
	//}
	msgs, err := repository.PullHistoryMsg(groupId, cursor)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
