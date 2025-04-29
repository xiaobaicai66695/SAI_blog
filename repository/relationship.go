package repository

type Relaationship struct {
	Uid int64 `gorm:"uid" json:"uid"`
	Fid int64 `gorm:"fid" json:"fid"`
}

func (Relaationship) TableName() string {
	return "relationship"
}

func SaveRelation(uid int64, fid int64) error {
	err := db.Create(&Relaationship{
		Uid: uid,
		Fid: fid,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteRelation(uid int64, fid int64) error {
	err := db.Where("uid = ? AND fid = ?", uid, fid).Delete(&Relaationship{}).Error
	if err != nil {
		return err
	}
	return nil
}

func IsRelationExist(uid int64, fid int64) (bool, error) {
	var count int64
	err := db.Model(&Relaationship{}).Where("uid = ? AND fid = ?", uid, fid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
