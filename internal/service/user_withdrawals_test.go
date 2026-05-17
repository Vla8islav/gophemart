package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMetricsService_GetUserWithdrawals(t *testing.T) {
	t.Parallel()

	processedAt := time.Date(2026, 5, 16, 10, 30, 0, 0,
		time.FixedZone("UTC+3", 3*60*60))

	tests := []struct {
		name            string
		userID          int64
		withdrawals     []domain.Withdrawal
		repoErr         error
		wantWithdrawals []domain.Withdrawal
		wantErr         error
	}{
		{
			name:   "success",
			userID: 42,
			withdrawals: []domain.Withdrawal{
				{
					OrderNumber: "2377225624",
					Amount:      50000,
					ProcessedAt: processedAt,
				},
			},
			wantWithdrawals: []domain.Withdrawal{
				{
					OrderNumber: "2377225624",
					Amount:      50000,
					ProcessedAt: processedAt,
				},
			},
		},
		{
			name:            "no withdrawals",
			userID:          42,
			withdrawals:     []domain.Withdrawal{},
			wantWithdrawals: []domain.Withdrawal{},
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
				GetUserWithdrawals(gomock.Any(), tt.userID).
				Return(tt.withdrawals, tt.repoErr)

			withdrawals, err := service.GetUserWithdrawals(context.Background(), tt.userID)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
				require.Nil(t, withdrawals)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantWithdrawals, withdrawals)
		})
	}
}
