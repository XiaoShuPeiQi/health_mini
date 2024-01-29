package manage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/global"
	"main.go/model/common/enum"
	"main.go/model/common/request"
	"main.go/model/common/response"
	manageReq "main.go/model/manage/request"
	manageRes "main.go/model/manage/response"
	"strconv"
)

type ManageGoodsCategoryApi struct {
}

// CreateCategory 新建商品分类
func (g *ManageGoodsCategoryApi) CreateCategory(c *gin.Context) {
	//1.获取参数
	var category manageReq.MallGoodsCategoryReq
	_ = c.ShouldBindJSON(&category)
	//2.service处理
	if err := mallGoodsCategoryService.AddCategory(category); err != nil {
		global.GVA_LOG.Error("创建失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// UpdateCategory 修改商品分类信息
func (g *ManageGoodsCategoryApi) UpdateCategory(c *gin.Context) {
	//1.获取参数
	var category manageReq.MallGoodsCategoryReq
	_ = c.ShouldBindJSON(&category)
	//2.service处理
	if err := mallGoodsCategoryService.UpdateCategory(category); err != nil {
		global.GVA_LOG.Error("更新失败", zap.Error(err))
		response.FailWithMessage("更新失败，存在相同分类", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// GetCategoryList 获取商品分类
func (g *ManageGoodsCategoryApi) GetCategoryList(c *gin.Context) {
	//1.获取参数
	var req manageReq.SearchCategoryParams
	_ = c.ShouldBindQuery(&req)
	//2.service处理
	if err, list, total := mallGoodsCategoryService.SelectCategoryPage(req); err != nil {
		global.GVA_LOG.Error("获取失败！", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
	} else {
		response.OkWithData(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   req.PageNumber,
			PageSize:   req.PageSize,
			TotalPage:  int(total) / req.PageSize,
		}, c)
	}
}

// GetCategory 通过id获取分类数据
func (g *ManageGoodsCategoryApi) GetCategory(c *gin.Context) {
	//1.获取参数
	id, _ := strconv.Atoi(c.Param("id"))
	//2.service处理
	err, goodsCategory := mallGoodsCategoryService.SelectCategoryById(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败！", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
	} else {
		response.OkWithData(manageRes.GoodsCategoryResponse{GoodsCategory: goodsCategory}, c)
	}
}

// DelCategory 设置分类失效
func (g *ManageGoodsCategoryApi) DelCategory(c *gin.Context) {
	//1.获取参数
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	//2.service处理
	if err, _ := mallGoodsCategoryService.DeleteGoodsCategoriesByIds(ids); err != nil {
		global.GVA_LOG.Error("删除失败！", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}

}

// ListForSelect 用于三级分类联动效果制作
func (g *ManageGoodsCategoryApi) ListForSelect(c *gin.Context) {
	//1.获取参数：商品分类id
	id, _ := strconv.Atoi(c.Param("id"))

	//2.根据商品分类ID查询商品分类信息
	err, goodsCategory := mallGoodsCategoryService.SelectCategoryById(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败！", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
	}
	//3.获取商品分类级别
	level := goodsCategory.CategoryLevel

	if level == enum.LevelThree.Code() ||
		level == enum.Default.Code() {
		response.FailWithMessage("参数异常", c)
	}
	//5.初始化结果集
	categoryResult := make(map[string]interface{})
	//6.根据商品分类级别进行不同的处理逻辑
	if level == enum.LevelOne.Code() {
		// 如果是一级分类，查询其下属的二级和三级分类
		_, levelTwoList := mallGoodsCategoryService.SelectByLevelAndParentIdsAndNumber(id, enum.LevelTwo.Code())
		if levelTwoList != nil {
			_, levelThreeList := mallGoodsCategoryService.SelectByLevelAndParentIdsAndNumber(int(levelTwoList[0].CategoryId), enum.LevelThree.Code())
			categoryResult["secondLevelCategories"] = levelTwoList
			categoryResult["thirdLevelCategories"] = levelThreeList
		}
	}
	if level == enum.LevelTwo.Code() {
		// 如果是二级分类，查询其下属的三级分类
		_, levelThreeList := mallGoodsCategoryService.SelectByLevelAndParentIdsAndNumber(id, enum.LevelThree.Code())
		categoryResult["thirdLevelCategories"] = levelThreeList
	}
	//7.返回结果
	response.OkWithData(categoryResult, c)
}
