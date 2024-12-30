package service

import (
	"distributed-order-system-go/internal/models"
	"distributed-order-system-go/internal/requests"
	"distributed-order-system-go/internal/responses"
	"distributed-order-system-go/pkg/global"

	"github.com/bwmarrin/snowflake"
)

var OrderService = NewOrderService()

func NewOrderService() *orderService {

	node, _ := snowflake.NewNode(1)
	return &orderService{
		IdGenerater: node,
	}
}

type orderService struct {
	IdGenerater *snowflake.Node
}

func (o *orderService) Create(orderDto *requests.CreateOrderDto) error {

	//redis 去重判断

	//验证商品、价格信息

	//分布式锁，减库存、生成订单

	order := mapToOrder(orderDto)

	if err := global.App.DB.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (o *orderService) List() {

}

func (o *orderService) GetOrderId() responses.OrderIdDto {
	return responses.OrderIdDto{
		ID: o.IdGenerater.Generate().Int64(),
	}
}

func mapToOrder(dto *requests.CreateOrderDto) *models.Order {

	orderItems := make([]models.OrderItem, 0, len(dto.Items))

	for _, item := range dto.Items {
		orderItems = append(orderItems, models.OrderItem{
			Amount:   item.Amount,
			Quantity: item.Quantity,
			Product:  item.Product,
			Price:    item.Price,
			OrderId:  dto.OrderId,
		})
	}

	return &models.Order{
		Amount:   dto.Amount,
		Customer: dto.Customer,
		Items:    orderItems,
	}
}
