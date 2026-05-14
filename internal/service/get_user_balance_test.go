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

func TestMetricsService_GetUserBalance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		userID      int64
		balance     *domain.UserBalance
		repoErr     error
		wantBalance *domain.UserBalance
		wantErr     error
	}{
		{
			name:   "success",
			userID: 42,
			balance: &domain.UserBalance{
				Current:   212452,
				Withdrawn: 1000,
			},
			wantBalance: &domain.UserBalance{
				Current:   212452,
				Withdrawn: 1000,
			},
		},
		{
			name:    "repository error",
			userID:  42,
			repoErr: errors.New("db is unavailable"),
			wantErr: errors.New("db is unavailable"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repository := mocks.NewMockGophermartRepository(ctrl)
			service := gophermartService{repository: repository}

			repository.EXPECT().
				GetUserBalance(gomock.Any(), tt.userID).
				Return(tt.balance, tt.repoErr)

			balance, err := service.GetUserBalance(context.Background(), tt.userID)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
				require.Nil(t, balance)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantBalance, balance)
		})
	}
}
