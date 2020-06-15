package main

import (
	"context"
	"fmt"
	"sync"

	orderServiceProtoGo "github.com/kandevarg/deliveryapp.orderservice/proto/protoGo"
	stockServiceProtoGo "github.com/kandevarg/deliveryapp.stockservice/proto/protoGo"
	goMicro "github.com/micro/go-micro"
)

type repository interface {
	createOrder(*orderServiceProtoGo.Order) (*orderServiceProtoGo.Order, error)
	getAllOrders() []*orderServiceProtoGo.Order
}
type dummyStore struct {
	mu     sync.RWMutex //semaphore
	orders []*orderServiceProtoGo.Order
}

func (ds *dummyStore) createOrder(order *orderServiceProtoGo.Order) (*orderServiceProtoGo.Order, error) {
	ds.mu.Lock()
	updated := append(ds.orders, order)
	ds.orders = updated
	ds.mu.Unlock()
	return order, nil
}

func (ds *dummyStore) getAllOrders() []*orderServiceProtoGo.Order {
	return ds.orders
}

type service struct {
	repo repository
}

func (svc *service) GetAllOrders(ctx context.Context, in *orderServiceProtoGo.BlankRequest, out *orderServiceProtoGo.GetOrdersResponse) error {
	out.Orders = svc.repo.getAllOrders()
	return nil
}

func (svc *service) CreateOrder(ctx context.Context, in *orderServiceProtoGo.Order, out *orderServiceProtoGo.CreateOrderResponse) error {

	fmt.Printf("Creating a new Order!!!! \n")
	order, err := svc.repo.createOrder(in)

	if err != nil {
		return err
	}

	out.Created = true
	out.Order = order
	return nil
}

func main() {

	repo := &dummyStore{}

	microService := goMicro.NewService(
		goMicro.Name("deliveryapp.orderservice"),
	)
	// Init will parse the command line flags.
	microService.Init()

	stockServiceClient := stockServiceProtoGo.NewStockServiceClient("deliveryapp.stockservice", microService.Client())

	greetings, err := stockServiceClient.Ping(context.Background(), &stockServiceProtoGo.PingRequest{
		CallerName: "OrderService",
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(greetings)

	// Register handler
	orderServiceProtoGo.RegisterOrderServiceHandler(microService.Server(), &service{repo})

	if err := microService.Run(); err != nil {
		fmt.Println(err)
	}
}
