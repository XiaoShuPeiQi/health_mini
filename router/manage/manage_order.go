package manage

import (
	"github.com/gin-gonic/gin"
	v1 "main.go/api/v1"
	"main.go/middleware"
)

type ManageOrderRouter struct {
}

func (r *ManageOrderRouter) InitManageOrderRouter(Router *gin.RouterGroup) {
	mallOrderRouter := Router.Group("v1").Use(middleware.AdminJWTAuth())
	var mallOrderApi = v1.ApiGroupApp.ManageApiGroup.ManageOrderApi
	{
		mallOrderRouter.PUT("orders/checkDone", mallOrderApi.CheckDoneOrder) // 发货(订单已支付完成)(code:2)
		mallOrderRouter.PUT("orders/checkOut", mallOrderApi.CheckOutOrder)   // 出库(code:3)
		mallOrderRouter.PUT("orders/close", mallOrderApi.CloseOrder)         // 商家关闭订单(code:-3)
		mallOrderRouter.GET("orders/:orderId", mallOrderApi.FindMallOrder)   // 根据ID获取订单信息
		mallOrderRouter.GET("orders", mallOrderApi.GetMallOrderList)         // 获取订单列表
	}
}
