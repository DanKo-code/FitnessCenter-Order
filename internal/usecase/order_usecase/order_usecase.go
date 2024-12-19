package order_usecase

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Order/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Order/internal/models"
	"github.com/DanKo-code/FitnessCenter-Order/internal/repository"
	abonementGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.abonement"
	orderGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.order"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/google/uuid"
	"time"
)

type OrderUseCase struct {
	orderRepo       repository.OrderRepository
	abonementClient *abonementGRPC.AbonementClient
	serviceClient   *serviceGRPC.ServiceClient
	userClient      *userGRPC.UserClient
}

func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	abonementClient *abonementGRPC.AbonementClient,
	userClient *userGRPC.UserClient,
	serviceClient *serviceGRPC.ServiceClient,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:       orderRepo,
		abonementClient: abonementClient,
		userClient:      userClient,
		serviceClient:   serviceClient,
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
		Status:      "Valid",
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}

	err = o.orderRepo.CreateCoachAbonement(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OrderUseCase) GetUserOrders(ctx context.Context, userId uuid.UUID) (*orderGRPC.GetUserOrdersResponse, error) {

	getUserByIdRequest := &userGRPC.GetUserByIdRequest{Id: userId.String()}

	_, err := (*o.userClient).GetUserById(ctx, getUserByIdRequest)
	if err != nil {
		return nil, err
	}

	orders, err := o.orderRepo.GetUserOrders(ctx, userId)
	if err != nil {
		return nil, err
	}

	var abonementIds []string
	for _, abonement := range orders {
		abonementIds = append(abonementIds, abonement.AbonementId.String())
	}

	abonements := &abonementGRPC.GetAbonementsByIdsResponse{}
	if len(abonementIds) > 0 {

		getAbonementsByIdsRequest := &abonementGRPC.GetAbonementsByIdsRequest{
			Ids: abonementIds,
		}

		abonements, err = (*o.abonementClient).GetAbonementsByIds(context.TODO(), getAbonementsByIdsRequest)
		if err != nil {
			return nil, err
		}
	}

	abonIdServicesMap := map[string][]*serviceGRPC.ServiceObject{}
	if len(abonementIds) > 0 {

		var abonIds []string
		for _, abon := range abonements.AbonementObjects {
			abonIds = append(abonIds, abon.Id)
		}

		getAbonementsServicesRequest := &serviceGRPC.GetAbonementsServicesRequest{
			AbonementIds: abonIds,
		}

		abonIdServices, err := (*o.serviceClient).GetAbonementsServices(ctx, getAbonementsServicesRequest)
		if err != nil {
			return nil, err
		}

		for _, object := range abonIdServices.AbonementIdsWithServices {
			abonIdServicesMap[object.AbonementId] = object.ServiceObjects
		}
	}

	abonIdOrder := map[string]*orderGRPC.OrderObject{}
	for _, order := range orders {
		orderObject := &orderGRPC.OrderObject{
			Id:          order.Id.String(),
			UserId:      order.UserId.String(),
			AbonementId: order.AbonementId.String(),
			Status:      order.Status,
			CreatedTime: order.CreatedTime.String(),
			UpdatedTime: order.UpdatedTime.String(),
		}

		abonIdOrder[orderObject.AbonementId] = orderObject
	}

	getUserOrdersResponse := orderGRPC.GetUserOrdersResponse{
		OrderObjectWithAbonementWithServices: nil,
	}

	for _, abonement := range abonements.AbonementObjects {

		orderObjectWithAbonement := &orderGRPC.OrderObjectWithAbonementWithServices{
			OrderObject:     abonIdOrder[abonement.Id],
			AbonementObject: abonement,
			ServiceObjects:  abonIdServicesMap[abonement.Id],
		}

		getUserOrdersResponse.OrderObjectWithAbonementWithServices = append(getUserOrdersResponse.OrderObjectWithAbonementWithServices, orderObjectWithAbonement)
	}

	return &getUserOrdersResponse, nil
}
