package service

import (
	"SAI_blog/repository"
	"fmt"
)

func Like(blogId int64, uid int64) error {
	isLike := repository.BlogIsLike(blogId, uid)
	if isLike {
		return fmt.Errorf("你已经点过赞了~~")
	}
	repository.Like(blogId, uid)
	return nil
}

func UnLike(blogId int64, uid int64) error {
	isLike := repository.BlogIsLike(blogId, uid)
	if !isLike {
		return fmt.Errorf("你还没点赞呢~~")
	}
	repository.UnLike(blogId, uid)
	return nil
}
