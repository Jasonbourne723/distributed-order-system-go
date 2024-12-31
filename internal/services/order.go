package service

import (
	"context"
	"distributed-order-system-go/internal/models"
	"distributed-order-system-go/internal/requests"
	"distributed-order-system-go/internal/responses"
	"distributed-order-system-go/pkg/global"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
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

func (o *orderService) Create(ctx context.Context, orderDto *requests.CreateOrderDto) error {

	//redis 去重判断
	key := "order:" + strconv.FormatInt(orderDto.OrderId, 10)
	if val, err := global.App.Redis.Incr(ctx, key).Result(); err != nil {
		return fmt.Errorf("Create Order fail %w ", err)
	} else if val > 1 {
		return errors.New("repeat order")
	}

	//验证商品、价格信息

	//分布式锁，事务处理 减库存、生成订单
	n := rand.Intn(20) + 1
	locker, err := global.App.DistributedLocker.WatchAndWaitForLock(orderDto.Items[0].Product + "-" + strconv.Itoa(n))
	if err != nil {
		return errors.New("get locker failed")
	}
	defer global.App.DistributedLocker.ReleaseLock(locker)

	return global.App.DB.Transaction(func(tx *gorm.DB) error {

		stock := "stock" + strconv.Itoa(n)
		result := tx.Model(&models.Inventory{}).Where("product = ?", orderDto.Items[0].Product).Where(stock+" > 0").UpdateColumn(stock, gorm.Expr(stock+" - ?", 1))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}
		order := mapToOrder(orderDto)
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		return nil
	})
}

func (o *orderService) List() {

}

func (o *orderService) GetOrderId() responses.OrderIdDto {
	return responses.OrderIdDto{
		ID: o.IdGenerater.Generate().Int64(),
	}
}

func mapToOrder(dto *requests.CreateOrderDto) models.Order {

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

	return models.Order{
		ID:       dto.OrderId,
		Amount:   dto.Amount,
		Customer: dto.Customer,
		Items:    orderItems,
		Created:  time.Now(),
	}
}
