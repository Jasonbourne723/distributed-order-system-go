from locust import HttpUser, task, between
import random
# 商品列表
product_ids = [
    "xz001", "xz002", "xz003", "xz004", "xz005"
]

class OrderUser(HttpUser):
    wait_time = between(1, 2)  # 每个请求之间等待 1 到 2 秒

    @task
    def create_order(self):
        # 获取订单号
        response = self.client.get("/order_id")
        if response.status_code == 200:
            order_id = response.json().get('id')

            product = random.choice(product_ids)
            # 下单
            self.client.post("/order", json={
                "order_id": order_id,
                "amount": 50,
                "customer": "vip1",
                "items": [
                    {"product": product, "quantity": 1, "price": 300, "amount": 300}
                ]
            })
