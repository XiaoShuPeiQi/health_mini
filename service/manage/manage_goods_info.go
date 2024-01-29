package manage

import (
	"errors"
	"gorm.io/gorm"
	"main.go/global"
	"main.go/model/common"
	"main.go/model/common/enum"
	"main.go/model/common/request"
	"main.go/model/manage"
	manageReq "main.go/model/manage/request"
	"main.go/utils"
	"strconv"
	"time"
)

type ManageGoodsInfoService struct {
}

// CreateMallGoodsInfo 创建MallGoodsInfo
func (m *ManageGoodsInfoService) CreateMallGoodsInfo(req manageReq.GoodsInfoAddParam) (err error) {
	//1.根据商品分类ID查询对应的商品分类信息
	var goodsCategory manage.MallGoodsCategory
	err = global.GVA_DB.Where("category_id=?  AND is_deleted=0", req.GoodsCategoryId).First(&goodsCategory).Error
	if goodsCategory.CategoryLevel != enum.LevelThree.Code() {
		return errors.New("分类数据异常")
	}
	//2.检查是否已存在相同的商品信息（根据商品名称和商品分类ID）
	if !errors.Is(global.GVA_DB.Where("goods_name=? AND goods_category_id=?", req.GoodsName, req.GoodsCategoryId).First(&manage.MallGoodsInfo{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("已存在相同的商品信息")
	}
	//3.一些类型的转化
	originalPrice, _ := strconv.Atoi(req.OriginalPrice)
	sellingPrice, _ := strconv.Atoi(req.SellingPrice)
	stockNum, _ := strconv.Atoi(req.StockNum)
	goodsSellStatus, _ := strconv.Atoi(req.GoodsSellStatus)
	//4.创建商品信息对象
	goodsInfo := manage.MallGoodsInfo{
		GoodsName:          req.GoodsName,
		GoodsIntro:         req.GoodsIntro,
		GoodsCategoryId:    req.GoodsCategoryId,
		GoodsCoverImg:      req.GoodsCoverImg,
		GoodsDetailContent: req.GoodsDetailContent,
		OriginalPrice:      originalPrice,
		SellingPrice:       sellingPrice,
		StockNum:           stockNum,
		Tag:                req.Tag,
		GoodsSellStatus:    goodsSellStatus,
		CreateTime:         common.JSONTime{Time: time.Now()},
		UpdateTime:         common.JSONTime{Time: time.Now()},
	}
	//5.校验商品信息对象的字段是否满足校验规则
	if err = utils.Verify(goodsInfo, utils.GoodsAddParamVerify); err != nil {
		return errors.New(err.Error())
	}
	//6.将商品信息对象插入数据库
	err = global.GVA_DB.Create(&goodsInfo).Error
	return err
}

// DeleteMallGoodsInfo 删除MallGoodsInfo记录
func (m *ManageGoodsInfoService) DeleteMallGoodsInfo(mallGoodsInfo manage.MallGoodsInfo) (err error) {
	err = global.GVA_DB.Delete(&mallGoodsInfo).Error
	return err
}

// ChangeMallGoodsInfoByIds 上下架
func (m *ManageGoodsInfoService) ChangeMallGoodsInfoByIds(ids request.IdsReq, sellStatus string) (err error) {
	//1.转化参数
	intSellStatus, _ := strconv.Atoi(sellStatus)
	//2.更新字段为0时，不能直接UpdateColumns,应该使用update
	err = global.GVA_DB.Model(&manage.MallGoodsInfo{}).Where("goods_id in ?", ids.Ids).Update("goods_sell_status", intSellStatus).Error
	return err
}

// UpdateMallGoodsInfo 更新MallGoodsInfo记录
func (m *ManageGoodsInfoService) UpdateMallGoodsInfo(req manageReq.GoodsInfoUpdateParam) (err error) {
	//1.数据转换
	goodsId, _ := strconv.Atoi(req.GoodsId)
	originalPrice, _ := strconv.Atoi(req.OriginalPrice)
	stockNum, _ := strconv.Atoi(req.StockNum)
	//2.定义goodsInfo变量
	goodsInfo := manage.MallGoodsInfo{
		GoodsId:            goodsId,
		GoodsName:          req.GoodsName,
		GoodsIntro:         req.GoodsIntro,
		GoodsCategoryId:    req.GoodsCategoryId,
		GoodsCoverImg:      req.GoodsCoverImg,
		GoodsDetailContent: req.GoodsDetailContent,
		OriginalPrice:      originalPrice,
		SellingPrice:       req.SellingPrice,
		StockNum:           stockNum,
		Tag:                req.Tag,
		GoodsSellStatus:    req.GoodsSellStatus,
		UpdateTime:         common.JSONTime{Time: time.Now()},
	}
	//3.检验字段合法性
	if err = utils.Verify(goodsInfo, utils.GoodsAddParamVerify); err != nil {
		return errors.New(err.Error())
	}
	//4.更新到数据库
	err = global.GVA_DB.Where("goods_id=?", goodsInfo.GoodsId).Updates(&goodsInfo).Error
	return err
}

// GetMallGoodsInfo 根据id获取第一条MallGoodsInfo记录
func (m *ManageGoodsInfoService) GetMallGoodsInfo(id int) (err error, mallGoodsInfo manage.MallGoodsInfo) {
	err = global.GVA_DB.Where("goods_id = ?", id).First(&mallGoodsInfo).Error
	return
}

// GetMallGoodsInfoInfoList 分页获取MallGoodsInfo记录
func (m *ManageGoodsInfoService) GetMallGoodsInfoInfoList(info manageReq.MallGoodsInfoSearch, goodsName string, goodsSellStatus string) (err error, list interface{}, total int64) {
	//1.指定参数
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	//2.创建db
	db := global.GVA_DB.Model(&manage.MallGoodsInfo{})
	//3.初始化结果集
	var mallGoodsInfos []manage.MallGoodsInfo
	//4.搜索list总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	// 5. 根据搜索条件修改db对象，构建查询
	if goodsName != "" {
		db.Where("goods_name =?", goodsName)
	}
	if goodsSellStatus != "" {
		db.Where("goods_sell_status =?", goodsSellStatus)
	}
	// 6. 根据分页参数、排序条件执行查询，获取商品信息列表
	err = db.Limit(limit).Offset(offset).Order("goods_id desc").Find(&mallGoodsInfos).Error
	return err, mallGoodsInfos, total
}
