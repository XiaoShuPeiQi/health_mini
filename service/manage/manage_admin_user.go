package manage

import (
	"errors"
	"gorm.io/gorm"
	"main.go/global"
	"main.go/model/manage"
	manageReq "main.go/model/manage/request"
	"main.go/utils"
	"strconv"
	"strings"
	"time"
)

type ManageAdminUserService struct {
}

// CreateMallAdminUser 创建MallAdminUser记录
func (m *ManageAdminUserService) CreateMallAdminUser(mallAdminUser manage.MallAdminUser) (err error) {
	//1.判断是否存在相同用户名
	//error.Is(a,b)用于判断a,b两个错误是否相同
	if !errors.Is(global.GVA_DB.Where("login_user_name = ?", mallAdminUser.LoginUserName).First(&manage.MallAdminUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同用户名")
	}
	//2.入库
	err = global.GVA_DB.Create(&mallAdminUser).Error
	return err
}

// UpdateMallAdminName 更新MallAdminUser昵称
func (m *ManageAdminUserService) UpdateMallAdminName(token string, req manageReq.MallUpdateNameParam) (err error) {
	// 声明并初始化 adminUserToken 变量，用于存储从数据库中获取的管理员用户令牌信息
	var adminUserToken manage.MallAdminUserToken
	// 通过 token 查询数据库，获取管理员用户令牌信息
	err = global.GVA_DB.Where("token =? ", token).First(&adminUserToken).Error
	if err != nil {
		return errors.New("不存在的用户")
	}
	// 使用获取到的 adminUserToken 中的 AdminUserId 字段值作为条件，更新数据库中对应的管理员用户的用户名和昵称
	err = global.GVA_DB.Where("admin_user_id = ?", adminUserToken.AdminUserId).Updates(&manage.MallAdminUser{
		LoginUserName: req.LoginUserName,
		NickName:      req.NickName,
	}).Error
	return err
}

func (m *ManageAdminUserService) UpdateMallAdminPassWord(token string, req manageReq.MallUpdatePasswordParam) (err error) {
	//1.验证token
	var adminUserToken manage.MallAdminUserToken
	err = global.GVA_DB.Where("token =? ", token).First(&adminUserToken).Error
	if err != nil {
		return errors.New("用户未登录")
	}
	//2.验证用户名
	var adminUser manage.MallAdminUser
	err = global.GVA_DB.Where("admin_user_id =?", adminUserToken.AdminUserId).First(&adminUser).Error
	if err != nil {
		return errors.New("不存在的用户")
	}
	//3.验证原密码
	if adminUser.LoginPassword != req.OriginalPassword {
		return errors.New("原密码不正确")
	}
	adminUser.LoginPassword = req.NewPassword
	//4.修改密码
	err = global.GVA_DB.Where("admin_user_id=?", adminUser.AdminUserId).Updates(&adminUser).Error
	return
}

// GetMallAdminUser 根据id获取MallAdminUser记录
func (m *ManageAdminUserService) GetMallAdminUser(token string) (err error, mallAdminUser manage.MallAdminUser) {
	//1.验证当前token合法性
	var adminToken manage.MallAdminUserToken
	if errors.Is(global.GVA_DB.Where("token =?", token).First(&adminToken).Error, gorm.ErrRecordNotFound) {
		return errors.New("不存在的用户"), mallAdminUser
	}
	//2.根据id查询user信息
	err = global.GVA_DB.Where("admin_user_id = ?", adminToken.AdminUserId).First(&mallAdminUser).Error
	return err, mallAdminUser
}

// AdminLogin 管理员登陆
func (m *ManageAdminUserService) AdminLogin(params manageReq.MallAdminLoginParam) (err error, mallAdminUser manage.MallAdminUser, adminToken manage.MallAdminUserToken) {
	//1.根据用户名密码查询登录信息，并保存到mallAdminUser{}中
	err = global.GVA_DB.Where("login_user_name=? AND login_password=?", params.UserName, params.PasswordMd5).First(&mallAdminUser).Error
	//2.判断mallAdminUser{}是否为空，作为查询成功与否的标志
	if mallAdminUser != (manage.MallAdminUser{}) {
		//(1)生成新token
		token := getNewToken(time.Now().UnixNano()/1e6, int(mallAdminUser.AdminUserId))
		//(2)查询旧token
		global.GVA_DB.Where("admin_user_id", mallAdminUser.AdminUserId).First(&adminToken)

		//(3)设置新增时间
		nowDate := time.Now()
		expireTime, _ := time.ParseDuration("48h")
		expireDate := nowDate.Add(expireTime)

		if adminToken == (manage.MallAdminUserToken{}) {
			//没有查询到adminToken，新增
			adminToken.AdminUserId = mallAdminUser.AdminUserId
			adminToken.Token = token
			adminToken.UpdateTime = nowDate
			adminToken.ExpireTime = expireDate
			if err = global.GVA_DB.Create(&adminToken).Error; err != nil {
				return
			}
		} else {
			//查询到token，更新
			adminToken.Token = token
			adminToken.UpdateTime = nowDate
			adminToken.ExpireTime = expireDate
			if err = global.GVA_DB.Save(&adminToken).Error; err != nil {
				return
			}
		}
	}
	return err, mallAdminUser, adminToken

}

// getNewToken 将当前时间和用户id为基础，生成token返回
func getNewToken(timeInt int64, userId int) (token string) {
	var build strings.Builder
	build.WriteString(strconv.FormatInt(timeInt, 10))
	build.WriteString(strconv.Itoa(userId))
	build.WriteString(utils.GenValidateCode(6))
	return utils.MD5V([]byte(build.String()))
}
