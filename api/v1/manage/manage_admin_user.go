package manage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/global"
	"main.go/model/common/request"
	"main.go/model/common/response"
	"main.go/model/example"
	"main.go/model/manage"
	manageReq "main.go/model/manage/request"
	"main.go/utils"
	"strconv"
)

type ManageAdminUserApi struct {
}

// CreateAdminUser 创建AdminUser
func (m *ManageAdminUserApi) CreateAdminUser(c *gin.Context) {
	//1.获取参数
	var params manageReq.MallAdminParam
	_ = c.ShouldBindJSON(&params)

	//2.validator检验参数的合法性
	if err := utils.Verify(params, utils.AdminUserRegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//3.封装结构体，传递给服务层
	mallAdminUser := manage.MallAdminUser{
		LoginUserName: params.LoginUserName,
		NickName:      params.NickName,
		LoginPassword: utils.MD5V([]byte(params.LoginPassword)),
	}
	//4.服务层尝试创建用户
	if err := mallAdminUserService.CreateMallAdminUser(mallAdminUser); err != nil {
		global.GVA_LOG.Error("创建失败:", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// 修改密码
func (m *ManageAdminUserApi) UpdateAdminUserPassword(c *gin.Context) {
	//1.获取参数
	var req manageReq.MallUpdatePasswordParam
	_ = c.ShouldBindJSON(&req)

	//是否加密？ 等项目出来再决定.....
	//mallAdminUserName := manage.MallAdminUser{
	//	LoginPassword: utils.MD5V([]byte(req.LoginPassword)),
	//}

	//2.获取token
	userToken := c.GetHeader("token")
	//3.交给service更新密码
	if err := mallAdminUserService.UpdateMallAdminPassWord(userToken, req); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}

}

// 更新用户名
func (m *ManageAdminUserApi) UpdateAdminUserName(c *gin.Context) {
	//1.获取参数
	var req manageReq.MallUpdateNameParam
	_ = c.ShouldBindJSON(&req)
	//2.获得userToken
	userToken := c.GetHeader("token")
	//3.交给服务层更新
	if err := mallAdminUserService.UpdateMallAdminName(userToken, req); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// AdminUserProfile 用id查询AdminUser
func (m *ManageAdminUserApi) AdminUserProfile(c *gin.Context) {
	//1.获取token
	adminToken := c.GetHeader("token")
	//2.交给service
	if err, mallAdminUser := mallAdminUserService.GetMallAdminUser(adminToken); err != nil {
		global.GVA_LOG.Error("未查询到记录", zap.Error(err))
		response.FailWithMessage("未查询到记录", c)
	} else {
		mallAdminUser.LoginPassword = "******"
		response.OkWithData(mallAdminUser, c)
	}
}

// AdminLogin 管理员登陆
func (m *ManageAdminUserApi) AdminLogin(c *gin.Context) {
	//1.参数绑定
	var adminLoginParams manageReq.MallAdminLoginParam
	_ = c.ShouldBindJSON(&adminLoginParams)
	adminLoginParams.PasswordMd5 = utils.MD5V([]byte(adminLoginParams.PasswordMd5))
	//2.交给服务层进行登录
	if err, _, adminToken := mallAdminUserService.AdminLogin(adminLoginParams); err != nil {
		response.FailWithMessage("登陆失败", c)
	} else {
		response.OkWithData(adminToken.Token, c)
	}
}

// AdminLogout 登出
func (m *ManageAdminUserApi) AdminLogout(c *gin.Context) {
	//通过删除数据库中的token来实现退出
	token := c.GetHeader("token")
	if err := mallAdminUserTokenService.DeleteMallAdminUserToken(token); err != nil {
		response.FailWithMessage("登出失败", c)
	} else {
		response.OkWithMessage("登出成功", c)
	}

}

// UserList 商城注册用户列表
func (m *ManageAdminUserApi) UserList(c *gin.Context) {
	//1.获取参数
	var pageInfo manageReq.MallUserSearch
	_ = c.ShouldBindQuery(&pageInfo)
	//2.service查询
	if err, list, total := mallUserService.GetMallUserInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// LockUser 用户禁用与解除禁用(0-未锁定 1-已锁定)
func (m *ManageAdminUserApi) LockUser(c *gin.Context) {
	//1.获取参数，从url里面获取参数
	lockStatus, _ := strconv.Atoi(c.Param("lockStatus"))
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)

	//2.service处理
	if err := mallUserService.LockUser(IDS, lockStatus); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// UploadFile 上传单图
// 此处上传图片的功能可用，但是原前端项目的图片链接为服务器地址，如需要显示图片，需要修改前端指向的图片链接
func (m *ManageAdminUserApi) UploadFile(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	//1.获取查询参数 noSave，用于指示是否保存文件，默认为 "0"
	noSave := c.DefaultQuery("noSave", "0")
	//2.获取需要上传的文件
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	//3.service文件上传
	//文件上传后：拿到文件路径
	err, file = fileUploadAndDownloadService.UploadFile(header, noSave)
	if err != nil {
		global.GVA_LOG.Error("修改数据库链接失败!", zap.Error(err))
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	//这里直接使用本地的url,后续项目开发时需要通过gorm将文件信息保存再数据库中....
	response.OkWithData("http://localhost:8888/"+file.Url, c)
}
