package manage

import (
	"github.com/gin-gonic/gin"
	v1 "main.go/api/v1"
)

type ManageGoodsInfoRouter struct {
}

func (m *ManageGoodsInfoRouter) InitManageGoodsInfoRouter(Router *gin.RouterGroup) {
	mallGoodsInfoRouter := Router.Group("v1")
	var mallGoodsInfoApi = v1.ApiGroupApp.ManageApiGroup.ManageGoodsInfoApi
	{

		mallGoodsInfoRouter.POST("goods", mallGoodsInfoApi.CreateGoodsInfo)                    // 添加商品信息
		mallGoodsInfoRouter.DELETE("deleteMallGoodsInfo", mallGoodsInfoApi.DeleteGoodsInfo)    // 删除商品信息
		mallGoodsInfoRouter.PUT("goods/status/:status", mallGoodsInfoApi.ChangeGoodsInfoByIds) // 修改商品状态(上下架)
		mallGoodsInfoRouter.PUT("goods", mallGoodsInfoApi.UpdateGoodsInfo)                     // 更新商品信息
		mallGoodsInfoRouter.GET("goods/:id", mallGoodsInfoApi.FindGoodsInfo)                   // 根据ID获取商品信息，包括分级分类信息
		mallGoodsInfoRouter.GET("goods/list", mallGoodsInfoApi.GetGoodsInfoList)               // 获取商品信息列表
	}
}
