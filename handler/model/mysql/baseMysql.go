//Package mysql
// @Author binwu.zhang 2022/1/25 8:33 下午
// @Description:
package mysql

import (
	"aig-tech-okr/libs/db"
	"gorm.io/gorm"
)

const (
	MysqlCluseterMaster = "master"
	MysqlCluseterSlave  = "slave"
)

type BaseMysql struct {
	Table     string
	ConnectDb *gorm.DB
	Cluster   string
}

func (i *BaseMysql) Connect() {

	if i.Cluster == "" {
		i.Cluster = MysqlCluseterMaster
	}

	if i.ConnectDb == nil {
		i.ConnectDb = db.GetMysqlConn(i.Cluster)
	}
	i.ConnectDb = i.ConnectDb.Table(i.Table)
}

//Insert
//  @date: 2022-02-22 13:56:32
//  @Description: 单条新增
//  @receiver i *BaseMysql
//  @param info *AIGDepartmentInfo
//  @return error
func (i *BaseMysql) Insert(info interface{}) error {
	return i.ConnectDb.Create(info).Error
}

//Inserts
//  @date: 2022-02-22 13:56:30
//  @Description: 批量新增
//  @receiver i *BaseMysql
//  @param info *[]AIGDepartmentInfo
//  @return error
func (i *BaseMysql) Inserts(list interface{}) error {
	return i.ConnectDb.Create(list).Error
}

//SaveById
//  @date: 2022-02-22 14:03:51
//  @Description: 根据主键Id更新，id=0则新增
//  @receiver i *BaseMysql
//  @param info interface{}
//  @return err error
func (i *BaseMysql) SaveById(info interface{}) (err error) {
	return i.ConnectDb.Save(info).Error
}

//UpdatesFieldsById
//  @date: 2022-02-22 14:08:27
//  @Description: 更新多个字段
//  @receiver i *BaseMysql
//  @param id uint
//  @param updateData map[string]interface{}
//  @return error
func (i *BaseMysql) UpdatesFieldsById(id uint, updateData map[string]interface{}) error {

	if id == 0 {
		return nil
	}
	return i.ConnectDb.Where("id=?", id).Updates(updateData).Error
}

//GetInfoById
//  @date: 2022-02-22 14:11:38
//  @Description: 根据Id获取详情
//  @receiver i *BaseMysql
//  @param id uint
//  @param info interface{}
//  @return err error
func (i *BaseMysql) GetInfoById(id uint, info interface{}) (err error) {

	if id == 0 {
		err = gorm.ErrRecordNotFound
		return
	}

	err = i.ConnectDb.Where("id=?", id).First(&info).Error
	return
}
