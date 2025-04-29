package service

import (
	"SAI_blog/repository"
	"fmt"
	"sync"
)

var mu sync.Mutex

func Follow(uid int64, fid int64) error {
	mu.Lock()
	defer mu.Unlock()

	exists, err := repository.IsRelationExist(uid, fid)
	if err != nil || exists {
		return fmt.Errorf("已关注")
	}
	err = repository.SaveRelation(uid, fid)
	repository.AddFollowBlog(uid, fid)
	if err != nil {
		return fmt.Errorf("关注失败")
	}
	return nil
}

func UnFollow(uid int64, fid int64) error {
	mu.Lock()
	defer mu.Unlock()

	exists, err := repository.IsRelationExist(uid, fid)
	if err != nil || !exists {
		return fmt.Errorf("未关注")
	}
	err = repository.DeleteRelation(uid, fid)
	repository.DelFollowBlog(uid, fid)
	if err != nil {
		return fmt.Errorf("取关失败")
	}
	return nil
}
