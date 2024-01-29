package initialize

import (
	"github.com/gin-gonic/gin"
	"main.go/global"
	"main.go/middleware"
	"main.go/router"
	"net/http"
)

func Routers() *gin.Engine {
	var Router = gin.Default()
	Router.StaticFS(global.GVA_CONFIG.Local.Path, http.Dir(global.GVA_CONFIG.Local.Path)) // 为用户头像和文件提供静态地址
	//Router.Use(middleware.LoadTls())  // 打开就能玩https了
	global.GVA_LOG.Info("use middleware logger")
	// 跨域
	Router.Use(middleware.Cors()) // 如需跨域可以打开
	global.GVA_LOG.Info("use middleware cors")
	// 方便统一添加路由组前缀 多服务器上线使用
	//商城后管路由
	manageRouter := router.RouterGroupApp.Manage
	ManageGroup := Router.Group("manage-api")
	PublicGroup := Router.Group("")

	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}
	{
		//商城后管路由初始化
		manageRouter.InitManageAdminUserRouter(ManageGroup)     // 管理员注册，登录，更新，上传图片等功能
		manageRouter.InitManageGoodsCategoryRouter(ManageGroup) //商品分类信息相关操作
		manageRouter.InitManageGoodsInfoRouter(ManageGroup)     //每个商品信息的操作：添加，删除，修改状态，更新信息，查询(byId,所有)
		manageRouter.InitManageCarouselRouter(ManageGroup)      //轮播图操作：添加，删除，更新，查询(byId，所有list)
		manageRouter.InitManageIndexConfigRouter(ManageGroup)   //首页配置项操作：添加，删除，更新，查询(byId，所有list)
		manageRouter.InitManageOrderRouter(ManageGroup)         //订单信息操作：发货，出库，商家关闭订单，查询(byId，所有list)
	}
	//商城前端路由
	mallRouter := router.RouterGroupApp.Mall
	MallGroup := Router.Group("api")
	{
		// 商城前端路由
		mallRouter.InitMallCarouselIndexRouter(MallGroup)
		mallRouter.InitMallGoodsInfoIndexRouter(MallGroup)
		mallRouter.InitMallGoodsCategoryIndexRouter(MallGroup)
		mallRouter.InitMallUserRouter(MallGroup)
		mallRouter.InitMallUserAddressRouter(MallGroup)
		mallRouter.InitMallShopCartRouter(MallGroup)
		mallRouter.InitMallOrderRouter(MallGroup)
	}
	global.GVA_LOG.Info("router register success")
	return Router
}
