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
	tx := repository.GetDB().Begin()
	exists, err := repository.IsRelationExist(uid, fid)
	if err != nil || exists {
		tx.Rollback()
		return fmt.Errorf("已关注")
	}
	err = repository.SaveRelation(tx, uid, fid)
	repository.AddFollowBlog(tx, uid, fid)
	user, err1 := repository.QueryUserById(uid)
	followerUser, err2 := repository.QueryUserById(fid)
	if err != nil || err1 != nil || err2 != nil {
		tx.Rollback()
		return fmt.Errorf("关注失败")
	}
	user.Follower += 1
	followerUser.Followed += 1
	err1 = repository.UpdateUser(tx, user)
	err2 = repository.UpdateUser(tx, followerUser)
	if err1 != nil || err2 != nil {
		tx.Rollback()
		return fmt.Errorf("关注后更新数据库时失败")
	}
	tx.Commit()
	return nil
}

func UnFollow(uid int64, fid int64) error {
	mu.Lock()
	defer mu.Unlock()
	tx := repository.GetDB().Begin()
	exists, err := repository.IsRelationExist(uid, fid)
	if err != nil || !exists {
		tx.Rollback()
		return fmt.Errorf("未关注")
	}
	err = repository.DeleteRelation(uid, fid)
	repository.DelFollowBlog(tx, fid)
	user, err1 := repository.QueryUserById(uid)
	followerUser, err2 := repository.QueryUserById(fid)
	if err != nil || err1 != nil || err2 != nil {
		tx.Rollback()
		return fmt.Errorf("取关失败")
	}
	user.Follower -= 1
	followerUser.Followed -= 1
	err1 = repository.UpdateUser(tx, user)
	err2 = repository.UpdateUser(tx, followerUser)
	if err1 != nil || err2 != nil {
		tx.Rollback()
		return fmt.Errorf("	取关后更新数据库时失败")
	}
	tx.Commit()
	return nil
}
