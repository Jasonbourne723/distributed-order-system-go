# distributed-order-system-go

这是一个基于 Go 语言开发的分布式订单系统。该项目提供了分布式锁和性能测试功能。

## 主要功能点

- 实现分布式锁机制,确保订单处理的一致性和并发性
- 提供性能测试工具,帮助开发者评估系统的性能表现
  
## 技术栈

- Gin
- GORM
- Zookeeper 分布式锁
- Redis

## 代码示例

### zookeeper 分布式锁的实现

```go
var (
	zkServers = []string{"127.0.0.1:2181"} // Zookeeper服务器地址
	lockPath  = "/lock"                    // 锁的路径
)

type ZookeeperLock struct {
	conn *zk.Conn
}

func NewZookeeperLock() *ZookeeperLock {
	conn, _, err := zk.Connect(zkServers, time.Second*35)
	if err != nil {
		fmt.Println(fmt.Errorf("zookeeper connect failed,%w", err))
	}
	return &ZookeeperLock{conn: conn}
}

// Watch and listen for the previous node in the sequence to be deleted.
func (zl *ZookeeperLock) WatchAndWaitForLock(name string) (string, error) {
	lockNode := lockPath + "/" + name + "-"
	nodeName, err := zl.conn.Create(lockNode, nil, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		return "", fmt.Errorf("unable to create lock node: %v", err)
	}
	for {
		children, _, err := zl.conn.Children(lockPath)
		if err != nil {
			return "", fmt.Errorf("unable to list lock children: %v", err)
		}

		// Sort nodes to determine the predecessor
		var previousNode string
		var currentNodeIdx int

		for idx, child := range children {
			if child == nodeName {
				currentNodeIdx = idx
				break
			}
		}

		// Watch for the predecessor node
		if currentNodeIdx > 0 {
			previousNode = children[currentNodeIdx-1]
		} else {
			// If there's no predecessor, it means this node is the first one, so no need to watch
			return nodeName, nil
		}

		// Watch the previous node for changes
		_, _, ch, err := zl.conn.ExistsW(lockPath + "/" + previousNode)
		if err != nil {
			return "", fmt.Errorf("unable to watch previous node: %v", err)
		}
		// Wait for the previous node to be deleted
		<-ch
	}

}

// Release the lock by deleting the lock node
func (zl *ZookeeperLock) ReleaseLock(nodeName string) error {
	err := zl.conn.Delete(nodeName, -1)
	if err != nil {
		return fmt.Errorf("unable to release lock: %v", err)
	}
	return nil
}
```
### 下单实现

```go
func (o *orderService) Create(ctx context.Context, orderDto *requests.CreateOrderDto) error {

	//redis 去重判断
	key := "order:" + strconv.FormatInt(orderDto.OrderId, 10)
	if val, err := global.App.Redis.Incr(ctx, key).Result(); err != nil {
		return fmt.Errorf("Create Order fail %w ", err)
	} else if val > 1 {
		return errors.New("repeat order")
	}
	global.App.Redis.Expire(ctx, key, time.Second*60)

	//验证商品、价格信息

	//分布式锁，事务处理 减库存、生成订单
	locker, err := global.App.DistributedLocker.WatchAndWaitForLock(orderDto.Items[0].Product)
	if err != nil {
		return errors.New("get locker failed")
	}
	defer global.App.DistributedLocker.ReleaseLock(locker)

	return global.App.DB.Transaction(func(tx *gorm.DB) error {

		result := tx.Model(&models.Inventory{}).Where("product = ?", orderDto.Items[0].Product).Where("stock > 0").UpdateColumn("stock", gorm.Expr("stock - ?", 1))
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
```

## 性能测试
![image](https://github.com/user-attachments/assets/fbc29a83-688c-4557-875e-b58a87b5ebcf)

- 使用 800 并发测试接口性能，qps 在 400 左右。
- 系统瓶颈在于数据库的事务处理
- 测试分段库存方案对 qps 没有任何提升，因为获取到分段的分布式锁之后，仍然需要等待数据库行锁。 

## 许可证

Apache-2.0许可证
