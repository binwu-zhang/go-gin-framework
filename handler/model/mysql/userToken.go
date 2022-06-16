package mysql

import (
	"gorm.io/gorm"
	"time"
)

type UserTokenModel struct {
	BaseMysql *BaseMysql
	ConnectDb *gorm.DB //事务链接db
}

type UserTokenInfo struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	UserToken  string    `gorm:"column:usertoken;type:varchar(100);NOT NULL" json:"usertoken"`                    // 用户token
	Openid     string    `gorm:"column:openid;type:varchar(100);NOT NULL" json:"openid"`                          // 用户id
	Platform   string    `gorm:"column:platform;type:varchar(50);NOT NULL" json:"platform"`                       //darwin代表Mac，windows_nt代表windows android  ios
	Expire     int64     `gorm:"column:expire;type:bigint(20);NOT NULL" json:"expire"`                            // 过期时间
	Created    time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"` //
	EmployeeId uint      `gorm:"column:employee_id;type:int(10) unsigned;NOT NULL" json:"employee_id"`            // aigEmployee
}

func (i *UserTokenModel) Base() *BaseMysql {
	i.BaseMysql = new(BaseMysql)
	i.BaseMysql.Table = i.TableName()
	i.BaseMysql.ConnectDb = i.ConnectDb
	i.BaseMysql.Connect()
	return i.BaseMysql
}

func (i *UserTokenModel) Connect() *gorm.DB {
	i.BaseMysql = new(BaseMysql)
	i.BaseMysql.Table = i.TableName()
	i.BaseMysql.ConnectDb = i.ConnectDb
	i.BaseMysql.Connect()
	return i.BaseMysql.ConnectDb
}

func (i *UserTokenModel) TableName() string {
	return "user_token"
}

func (i *UserTokenModel) GetInfoByOpenid(openid string, platform string) (info UserTokenInfo, err error) {
	err = i.Connect().Where("openid=? and platform=?", openid, platform).First(&info).Error
	return
}

func (i *UserTokenModel) GetInfoByToken(token string, platform string) (info UserTokenInfo, err error) {
	err = i.Connect().Where("usertoken=? and platform=?", token, platform).First(&info).Error
	return
}
