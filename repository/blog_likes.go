package repository

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type BlogLikes struct {
	BlogId int64 `gorm:"column:blog_id" json:"blog_id"`
	Lid    int64 `gorm:"column:lid" json:"lid"`
}

func QueryLikesById(blogId int64) []BlogLikes {
	var likes []BlogLikes
	db.Where("blog_id=?", blogId).Find(&likes)
	return likes
}

func BlogIsLike(blogId int64, lid int64) bool {
	key := fmt.Sprintf("blog:like:%d", blogId)
	ok, _ := rdb2.SIsMember(ctx, key, lid).Result()
	if ok {
		return true
	}
	var like BlogLikes
	db.Where("blog_id=? and lid = ?", blogId, lid).Find(&like)
	if like.Lid == lid {
		rdb2.SAdd(ctx, key, lid)
		return true
	}
	return false
}
func likeCountFunc(blogId int64) {
	keyCount := fmt.Sprintf("blog:likes:count:%d", blogId)
	db.Model(&Blog{}).Where("blog_id = ?", blogId).Find(&blog)
	rdb2.Set(ctx, keyCount, blog.Likes, -1)
}
func Like(blogId int64, lid int64) {
	key := fmt.Sprintf("blog:like:%d", blogId)
	rdb2.SAdd(ctx, key, lid)
	keyCount := fmt.Sprintf("blog:likes:count:%d", blogId)
	db.Model(&BlogLikes{}).Create(&BlogLikes{BlogId: blogId, Lid: lid})
	rdb2.IncrBy(ctx, keyCount, 1)
	return
}

func UnLike(blogId int64, lid int64) {
	key := fmt.Sprintf("blog:like:%d", blogId)
	keyCount := fmt.Sprintf("blog:likes:count:%d", blogId)
	rdb2.SRem(ctx, key, lid)
	rdb2.IncrBy(ctx, keyCount, -1)
	db.Model(&BlogLikes{}).Where("blog_id = ?", blogId).Delete(&BlogLikes{BlogId: blogId, Lid: lid})
	return
}

func UpdateLikeFromRedis() {
	var cursorCount uint64
	var likeCountKeys []string
	var err error
	for {
		likeCountKeys, cursorCount, err = rdb2.Scan(ctx, cursorCount, "blog:likes:count:*", 100).Result()
		if err != nil {
			break
		}
		for _, likeCountKey := range likeCountKeys {
			likeCountStr, _ := rdb2.Get(ctx, likeCountKey).Result()
			parts := strings.Split(likeCountKey, ":")
			countKey, _ := strconv.ParseInt(parts[3], 10, 64)
			likeCount, _ := strconv.ParseInt(likeCountStr, 10, 64)
			err = db.Model(&blog).Where("blog_id = ?", countKey).Update("likes", gorm.Expr("likes+?", likeCount)).Error
			rdb2.Del(ctx, likeCountKey)
			if err != nil {
				print(err)
			}
		}
		if cursorCount == 0 {
			break
		}
	}
}
