package manage

import (
	"errors"
	"github.com/jinzhu/copier"
	"main.go/global"
	"main.go/model/common"
	"main.go/model/common/enum"
	"main.go/model/common/request"
	"main.go/model/manage"
	manageRes "main.go/model/manage/response"
	"strconv"
	"time"
)

type ManageOrderService struct {
}

// CheckDone 修改订单状态为配货成功
func (m *ManageOrderService) CheckDone(ids request.IdsReq) (err error) {
	//1.根据id查询订单信息
	var orders []manage.MallOrder
	err = global.GVA_DB.Where("order_id in ?", ids.Ids).Find(&orders).Error
	var errorOrders string
	//2.遍历订单列表进行检查
	if len(orders) != 0 {
		for _, order := range orders {
			// 如果订单已被删除，则记录下订单号，跳过后续检查
			if order.IsDeleted == 1 {
				errorOrders = order.OrderNo + " "
				continue
			}
			// 如果订单状态不是已支付成功，则记录下订单号，表示无法执行出库操作
			if order.OrderStatus != enum.ORDER_PAID.Code() {
				errorOrders = order.OrderNo + " "
			}
		}
		// 如果未记录任何异常订单，则执行出库操作
		if errorOrders == "" {
			// 将指定订单ID列表对应的订单状态更新为已完成（OrderStatus=2），并更新更新时间
			if err = global.GVA_DB.Where("order_id in ?", ids.Ids).UpdateColumns(manage.MallOrder{OrderStatus: 2, UpdateTime: common.JSONTime{Time: time.Now()}}).Error; err != nil {
				return err
			}
		} else {
			// 如果存在异常订单，则返回错误信息，表示无法执行出库操作
			return errors.New("订单的状态不是支付成功无法执行出库操作")
		}
	}
	return
}

// CheckOut 出库
func (m *ManageOrderService) CheckOut(ids request.IdsReq) (err error) {
	var orders []manage.MallOrder
	err = global.GVA_DB.Where("order_id in ?", ids.Ids).Find(&orders).Error

	var errorOrders string
	if len(orders) != 0 {
		for _, order := range orders {
			if order.IsDeleted == 1 {
				errorOrders = order.OrderNo + " "
				continue
			}
			// 如果订单状态不是已支付成功或配货完成，则记录下订单号，表示无法执行出库操作
			if order.OrderStatus != enum.ORDER_PAID.Code() && order.OrderStatus != enum.ORDER_PACKAGED.Code() {
				errorOrders = order.OrderNo + " "
			}
		}
		if errorOrders == "" {
			// 将指定订单ID列表对应的订单状态更新为已出库（OrderStatus=3），并更新更新时间
			if err = global.GVA_DB.Where("order_id in ?", ids.Ids).UpdateColumns(manage.MallOrder{OrderStatus: 3, UpdateTime: common.JSONTime{Time: time.Now()}}).Error; err != nil {
				return err
			}
		} else {
			return errors.New("订单的状态不是支付成功或配货完成无法执行出库操作")
		}
	}
	return
}

// CloseOrder 商家关闭订单
func (m *ManageOrderService) CloseOrder(ids request.IdsReq) (err error) {
	var orders []manage.MallOrder
	err = global.GVA_DB.Where("order_id in ?", ids.Ids).Find(&orders).Error
	var errorOrders string
	if len(orders) != 0 {
		for _, order := range orders {
			if order.IsDeleted == 1 {
				errorOrders = order.OrderNo + " "
				continue
			}
			if order.OrderStatus == enum.ORDER_SUCCESS.Code() || order.OrderStatus < 0 {
				errorOrders = order.OrderNo + " "
			}
		}
		if errorOrders == "" {
			// 将指定订单ID列表对应的订单状态更新为已关闭（code:-3），并更新更新时间
			if err = global.GVA_DB.Where("order_id in ?", ids.Ids).UpdateColumns(manage.MallOrder{OrderStatus: enum.ORDER_CLOSED_BY_JUDGE.Code(), UpdateTime: common.JSONTime{Time: time.Now()}}).Error; err != nil {
				return err
			}
		} else {
			return errors.New("订单不能执行关闭操作")
		}
	}
	return
}

// GetMallOrder 根据id获取MallOrder记录
func (m *ManageOrderService) GetMallOrder(id string) (err error, newBeeMallOrderDetailVO manageRes.NewBeeMallOrderDetailVO) {
	// 1.查询指定订单ID对应的订单信息
	var newBeeMallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_id = ?", id).First(&newBeeMallOrder).Error; err != nil {
		return
	}

	// 查询指定订单ID对应的订单项列表信息
	var orderItems []manage.MallOrderItem
	if err = global.GVA_DB.Where("order_id = ?", newBeeMallOrder.OrderId).Find(&orderItems).Error; err != nil {
		return
	}
	//获取订单项数据
	if len(orderItems) > 0 {
		var newBeeMallOrderItemVOS []manageRes.NewBeeMallOrderItemVO
		// 将订单项数据复制到新的结构体中
		copier.Copy(&newBeeMallOrderItemVOS, &orderItems)
		// 将订单数据和订单项数据复制到新的结构体中
		copier.Copy(&newBeeMallOrderDetailVO, &newBeeMallOrder)

		// 获取订单状态和支付方式的字符串表示，并添加到新的结构体中
		_, OrderStatusStr := enum.GetNewBeeMallOrderStatusEnumByStatus(newBeeMallOrderDetailVO.OrderStatus)
		_, payTapStr := enum.GetNewBeeMallOrderStatusEnumByStatus(newBeeMallOrderDetailVO.PayType)

		newBeeMallOrderDetailVO.OrderStatusString = OrderStatusStr
		newBeeMallOrderDetailVO.PayTypeString = payTapStr
		newBeeMallOrderDetailVO.NewBeeMallOrderItemVOS = newBeeMallOrderItemVOS
	}
	return
}

// GetMallOrderInfoList 分页获取MallOrder记录
func (m *ManageOrderService) GetMallOrderInfoList(info request.PageInfo, orderNo string, orderStatus string) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	// 1.创建db
	db := global.GVA_DB.Model(&manage.MallOrder{})
	if orderNo != "" {
		db.Where("order_no", orderNo)
	}
	// 0.待支付 1.已支付 2.配货完成 3:出库成功 4.交易成功 -1.手动关闭 -2.超时关闭 -3.商家关闭
	if orderStatus != "" {
		status, _ := strconv.Atoi(orderStatus)
		db.Where("order_status", status)
	}
	var mallOrders []manage.MallOrder
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("update_time desc").Find(&mallOrders).Error
	return err, mallOrders, total
}
