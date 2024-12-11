package order_usecase

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Order/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Order/internal/models"
	"github.com/DanKo-code/FitnessCenter-Order/internal/repository"
	abonementGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.abonement"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/google/uuid"
	"time"
)

type OrderUseCase struct {
	orderRepo       repository.OrderRepository
	abonementClient *abonementGRPC.AbonementClient
	userClient      *userGRPC.UserClient
}

func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	abonementClient *abonementGRPC.AbonementClient,
	userClient *userGRPC.UserClient,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:       orderRepo,
		abonementClient: abonementClient,
		userClient:      userClient,
	}
}

func (o *OrderUseCase) CreateCoachAbonement(ctx context.Context, cmd *dtos.CreateOrderCommand) (*models.Order, error) {

	getUserByIdRequest := &userGRPC.GetUserByIdRequest{Id: cmd.UserId.String()}

	_, err := (*o.userClient).GetUserById(ctx, getUserByIdRequest)
	if err != nil {
		return nil, err
	}

	getAbonementByIdRequest := &abonementGRPC.GetAbonementByIdRequest{Id: cmd.AbonementId.String()}

	_, err = (*o.abonementClient).GetAbonementById(ctx, getAbonementByIdRequest)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		Id:          uuid.New(),
		AbonementId: cmd.AbonementId,
		UserId:      cmd.UserId,
		Status:      "In Process",
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}

	err = o.orderRepo.CreateCoachAbonement(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OrderUseCase) GetUserOrders(ctx context.Context, userId uuid.UUID) ([]*models.Order, error) {

	getUserByIdRequest := &userGRPC.GetUserByIdRequest{Id: userId.String()}

	_, err := (*o.userClient).GetUserById(ctx, getUserByIdRequest)
	if err != nil {
		return nil, err
	}

	abonements, err := o.orderRepo.GetUserOrders(ctx, userId)
	if err != nil {
		return nil, err
	}

	return abonements, nil
}
