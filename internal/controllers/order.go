package controllers

import (
	"distributed-order-system-go/internal/requests"
	"distributed-order-system-go/internal/responses"
	service "distributed-order-system-go/internal/services"

	"github.com/gin-gonic/gin"
)

var OrderApi = new(orderApi)

type orderApi struct{}

func (o *orderApi) Create(c *gin.Context) {
	var createOrderDto requests.CreateOrderDto
	if err := c.ShouldBindBodyWithJSON(&createOrderDto); err != nil {
		c.JSON(400, responses.ErrResult{
			Err: err.Error(),
		})
		return
	}

	if err := service.OrderService.Create(c, &createOrderDto); err != nil {
		c.JSON(500, responses.ErrResult{
			Err: err.Error(),
		})
		return
	}
	c.Status(200)
}

func (o *orderApi) GetOrderId(c *gin.Context) {
	orderId := service.OrderService.GetOrderId()
	c.JSON(200, orderId)
}
