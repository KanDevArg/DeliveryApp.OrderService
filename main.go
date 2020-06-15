package main

import (
	"context"
	"fmt"
	"sync"

	osModel "github.com/kandevarg/deliveryapp.orderservice/proto/protoGo"
	micro "github.com/micro/go-micro"
)

// const (
// 	port = ":50054"
// )

type repository interface {
	Create(*osModel.Order) (*osModel.Order, error)
	GetAll() []*osModel.Order
}

// DummyStore - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type DummyStore struct {
	mu     sync.RWMutex //semaphore
	orders []*osModel.Order
}

// Create a new order
func (ds *DummyStore) Create(order *osModel.Order) (*osModel.Order, error) {
	ds.mu.Lock()
	updated := append(ds.orders, order)
	ds.orders = updated
	ds.mu.Unlock()
	return order, nil
}

// GetAll consignments
func (ds *DummyStore) GetAll() []*osModel.Order {
	return ds.orders
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo repository
}

// GetOrders : Get all saved orders.
func (svc *service) GetOrders(ctx context.Context, in *osModel.GetRequest, out *osModel.GetOrdersResponse) error {
	out.Orders = svc.repo.GetAll()
	//return &osModel.Response{Orders: orders}, nil
	return nil
}

// CreateOrder - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (svc *service) CreateOrder(ctx context.Context, in *osModel.Order, out *osModel.CreateOrderResponse) error {

	fmt.Printf("Creating a new Order!!!! \n")

	// Save our order
	order, err := svc.repo.Create(in)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	out.Created = true
	out.Order = order
	return nil
}

func main() {
	repo := &DummyStore{}

	microService := micro.NewService(
		micro.Name("deliveryapp.productorderservice"),
	)
	// Init will parse the command line flags.
	microService.Init()

	// Register handler
	osModel.RegisterOrderServiceHandler(microService.Server(), &service{repo})

	// Run the servers
	if err := microService.Run(); err != nil {
		fmt.Println(err)
	}
}
