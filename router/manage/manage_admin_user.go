package manage

import (
	"github.com/gin-gonic/gin"
	v1 "main.go/api/v1"
	"main.go/middleware"
)

type ManageAdminUserRouter struct {
}

func (r *ManageAdminUserRouter) InitManageAdminUserRouter(Router *gin.RouterGroup) {
	mallAdminUserRouter := Router.Group("v1").Use(middleware.AdminJWTAuth())
	mallAdminUserWithoutRouter := Router.Group("v1")
	var mallAdminUserApi = v1.ApiGroupApp.ManageApiGroup.ManageAdminUserApi
	{

		mallAdminUserRouter.PUT("adminUser/name", mallAdminUserApi.UpdateAdminUserName)         // 更新MallAdminUser
		mallAdminUserRouter.PUT("adminUser/password", mallAdminUserApi.UpdateAdminUserPassword) // 修改密码
		mallAdminUserRouter.GET("users", mallAdminUserApi.UserList)                             // 获取管理员用户列表
		mallAdminUserRouter.PUT("users/:lockStatus", mallAdminUserApi.LockUser)                 // 修改锁状态
		mallAdminUserRouter.GET("adminUser/profile", mallAdminUserApi.AdminUserProfile)         // 根据id查询user信息
		mallAdminUserRouter.DELETE("logout", mallAdminUserApi.AdminLogout)                      // 退出登录
		mallAdminUserRouter.POST("upload/file", mallAdminUserApi.UploadFile)                    // 上传图片

	}
	{
		mallAdminUserWithoutRouter.POST("createMallAdminUser", mallAdminUserApi.CreateAdminUser) // 注册MallAdminUser
		mallAdminUserWithoutRouter.POST("adminUser/login", mallAdminUserApi.AdminLogin)          // 管理员登陆
	}
}
