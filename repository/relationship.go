package repository

import "gorm.io/gorm"

type Relationship struct {
	//关注者
	Uid int64 `gorm:"uid" json:"uid"`
	//被关注者
	Fid int64 `gorm:"fid" json:"fid"`
}

func (Relationship) TableName() string {
	return "relationship"
}

func SaveRelation(tx *gorm.DB, uid int64, fid int64) error {
	err := tx.Create(&Relationship{
		Uid: uid,
		Fid: fid,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteRelation(uid int64, fid int64) error {
	err := db.Where("uid = ? AND fid = ?", uid, fid).Delete(&Relationship{}).Error
	if err != nil {
		return err
	}
	return nil
}

func IsRelationExist(uid int64, fid int64) (bool, error) {
	var count int64
	err := db.Model(&Relationship{}).Where("uid = ? AND fid = ?", uid, fid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
