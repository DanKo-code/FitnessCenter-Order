package grpc

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Order/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Order/internal/usecase"
	orderProtobuf "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.order"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type OrdergRPC struct {
	orderProtobuf.UnimplementedOrderServer

	orderUseCase usecase.OrderUseCase
}

func Register(gRPC *grpc.Server, orderUseCase usecase.OrderUseCase) {
	orderProtobuf.RegisterOrderServer(gRPC, &OrdergRPC{orderUseCase: orderUseCase})
}

func (o OrdergRPC) CreateOrder(ctx context.Context, request *orderProtobuf.CreateOrderRequest) (*orderProtobuf.CreateOrderResponse, error) {

	cmd := &dtos.CreateOrderCommand{
		UserId:      uuid.MustParse(request.OrderDataForCreate.UserId),
		AbonementId: uuid.MustParse(request.OrderDataForCreate.AbonementId),
	}

	order, err := o.orderUseCase.CreateCoachAbonement(ctx, cmd)
	if err != nil {
		return nil, err
	}

	orderObject := &orderProtobuf.OrderObject{
		Id:          order.Id.String(),
		UserId:      order.UserId.String(),
		AbonementId: order.AbonementId.String(),
		Status:      order.Status,
		CreatedTime: order.CreatedTime.String(),
		UpdatedTime: order.UpdatedTime.String(),
	}

	response := &orderProtobuf.CreateOrderResponse{
		OrderObject: orderObject,
	}

	return response, nil
}

func (o OrdergRPC) GetUserOrders(ctx context.Context, request *orderProtobuf.GetUserOrdersRequest) (*orderProtobuf.GetUserOrdersResponse, error) {
	orders, err := o.orderUseCase.GetUserOrders(ctx, uuid.MustParse(request.UserId))
	if err != nil {
		return nil, err
	}

	var orderObjects []*orderProtobuf.OrderObject
	for _, order := range orders {

		orderObject := &orderProtobuf.OrderObject{
			Id:          order.Id.String(),
			UserId:      order.UserId.String(),
			AbonementId: order.AbonementId.String(),
			Status:      order.Status,
			CreatedTime: order.CreatedTime.String(),
			UpdatedTime: order.UpdatedTime.String(),
		}

		orderObjects = append(orderObjects, orderObject)
	}

	response := &orderProtobuf.GetUserOrdersResponse{
		UserObjects: orderObjects,
	}

	return response, nil
}
