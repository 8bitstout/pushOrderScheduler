# OrderPushScheduler
OrderPushScheduler is a service which receives new customer orders
and schedules a push notification for that customer's device to notify
that customer of when their order is ready to be picked up in a store.

ORDER_CREATION_SERVICE -> OrderPushScheduler -> PUSH_TO_DEVICE_SERVICE

## Architecture
OrderPushSceduler uses Protobuf for network messages and gRPC for message streaming.

## Service Definitions
### ScheduleOrderMessage
Receive an order and schedule a push notification using order details.
