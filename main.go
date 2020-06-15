package main

import (
	"context"
	"fmt"
	"sync"

	protoGo "github.com/kandevarg/deliveryapp.orderservice/proto/protoGo"
	goMicro "github.com/micro/go-micro"
)

type repository interface {
	createOrder(*protoGo.Order) (*protoGo.Order, error)
	getAllOrders() []*protoGo.Order
}
type dummyStore struct {
	mu     sync.RWMutex //semaphore
	orders []*protoGo.Order
}

func (ds *dummyStore) createOrder(order *protoGo.Order) (*protoGo.Order, error) {
	ds.mu.Lock()
	updated := append(ds.orders, order)
	ds.orders = updated
	ds.mu.Unlock()
	return order, nil
}

func (ds *dummyStore) getAllOrders() []*protoGo.Order {
	return ds.orders
}

type service struct {
	repo repository
}

func (svc *service) GetAllOrders(ctx context.Context, in *protoGo.BlankRequest, out *protoGo.GetOrdersResponse) error {
	out.Orders = svc.repo.getAllOrders()
	return nil
}

func (svc *service) CreateOrder(ctx context.Context, in *protoGo.Order, out *protoGo.CreateOrderResponse) error {

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

	// Register handler
	protoGo.RegisterOrderServiceHandler(microService.Server(), &service{repo})

	if err := microService.Run(); err != nil {
		fmt.Println(err)
	}
}
