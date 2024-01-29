package manage

import (
	"github.com/gin-gonic/gin"
	v1 "main.go/api/v1"
	"main.go/middleware"
)

type ManageGoodsCategoryRouter struct {
}

func (r *ManageGoodsCategoryRouter) InitManageGoodsCategoryRouter(Router *gin.RouterGroup) {
	goodsCategoryRouter := Router.Group("v1").Use(middleware.AdminJWTAuth())

	var goodsCategoryApi = v1.ApiGroupApp.ManageApiGroup.ManageGoodsCategoryApi
	{
		goodsCategoryRouter.POST("categories", goodsCategoryApi.CreateCategory)      //添加商品的分类信息
		goodsCategoryRouter.PUT("categories", goodsCategoryApi.UpdateCategory)       //修改商品的分类信息
		goodsCategoryRouter.GET("categories", goodsCategoryApi.GetCategoryList)      //获取商品分类列表
		goodsCategoryRouter.GET("categories/:id", goodsCategoryApi.GetCategory)      //通过id获取分类数据
		goodsCategoryRouter.DELETE("categories", goodsCategoryApi.DelCategory)       //通过分类id删除某种分类
		goodsCategoryRouter.GET("categories4Select", goodsCategoryApi.ListForSelect) //获取商品分类及下属的二级三级分类信息
	}
}
