package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGophermartService_WithdrawFromUserBalance(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := int64(123)
	request := domain.UserBalanceWithdraw{
		Order: "2377225624",
		Sum:   751,
	}

	repository := mocks.NewMockGophermartRepository(ctrl)
	repository.EXPECT().
		WithdrawFromUserBalance(gomock.Any(), userID, request).
		Return(nil)

	service := gophermartService{
		repository: repository,
	}

	err := service.WithdrawFromUserBalance(ctx, userID, request)
	require.NoError(t, err)
}

func TestGophermartService_WithdrawFromUserBalance_InvalidSum(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := mocks.NewMockGophermartRepository(ctrl)

	service := gophermartService{
		repository: repository,
	}

	err := service.WithdrawFromUserBalance(ctx, 123, domain.UserBalanceWithdraw{
		Order: "2377225624",
		Sum:   0,
	})
	require.ErrorIs(t, err, domain.ErrInvalidSum)
}

func TestGophermartService_WithdrawFromUserBalance_EmptyOrder(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := mocks.NewMockGophermartRepository(ctrl)

	service := gophermartService{
		repository: repository,
	}

	err := service.WithdrawFromUserBalance(ctx, 123, domain.UserBalanceWithdraw{
		Order: "",
		Sum:   751,
	})
	require.ErrorIs(t, err, domain.ErrOrderEmpty)
}

func TestGophermartService_WithdrawFromUserBalance_InvalidOrderNumber(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := mocks.NewMockGophermartRepository(ctrl)

	service := gophermartService{
		repository: repository,
	}

	err := service.WithdrawFromUserBalance(ctx, 123, domain.UserBalanceWithdraw{
		Order: "1234567890",
		Sum:   751,
	})
	require.ErrorIs(t, err, domain.ErrInvalidOrderNumber)
}

func TestGophermartService_WithdrawFromUserBalance_RepositoryError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := int64(123)
	request := domain.UserBalanceWithdraw{
		Order: "2377225624",
		Sum:   751,
	}
	repositoryErr := errors.New("repository error")

	repository := mocks.NewMockGophermartRepository(ctrl)
	repository.EXPECT().
		WithdrawFromUserBalance(gomock.Any(), userID, request).
		Return(repositoryErr)

	service := gophermartService{
		repository: repository,
	}

	err := service.WithdrawFromUserBalance(ctx, userID, request)
	require.ErrorIs(t, err, repositoryErr)
}
