package manage

import (
	"errors"
	"main.go/global"
	"main.go/model/common/request"
	"main.go/model/manage"
	manageReq "main.go/model/manage/request"
)

type ManageUserService struct {
}

// LockUser 修改用户状态
func (m *ManageUserService) LockUser(idReq request.IdsReq, lockStatus int) (err error) {
	//1.验证参数合法性
	if lockStatus != 0 && lockStatus != 1 {
		return errors.New("操作非法！")
	}
	//2.更新数据库
	//更新字段为0时，不能直接UpdateColumns方法，所以使用Update
	err = global.GVA_DB.Model(&manage.MallUser{}).Where("user_id in ?", idReq.Ids).Update("locked_flag", lockStatus).Error
	return err
}

// GetMallUserInfoList 分页获取商城注册用户列表
func (m *ManageUserService) GetMallUserInfoList(info manageReq.MallUserSearch) (err error, list interface{}, total int64) {
	//1.计算limit和偏移量offset
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	//2.创建db
	db := global.GVA_DB.Model(&manage.MallUser{})
	var mallUsers []manage.MallUser

	//3.获取用户总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	//4.获取列表，按照指定限制返回
	err = db.Limit(limit).Offset(offset).Order("create_time desc").Find(&mallUsers).Error
	return err, mallUsers, total
}
