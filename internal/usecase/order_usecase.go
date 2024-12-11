package usecase

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Order/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Order/internal/models"
	"github.com/google/uuid"
)

type OrderUseCase interface {
	CreateCoachAbonement(ctx context.Context, cmd *dtos.CreateOrderCommand) (*models.Order, error)

	GetUserOrders(ctx context.Context, userId uuid.UUID) ([]*models.Order, error)
}
