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

func TestGophermartService_CreateOrder(t *testing.T) {
	t.Parallel()

	const userID int64 = 42
	repoGenericErr := errors.New("database is unavailable")

	tests := []struct {
		name             string
		orderNumber      string
		repoErr          error
		wantErr          error
		wantRepoCall     bool
		wantRepoOrderNum string
	}{
		{
			name:             "success",
			orderNumber:      "79927398713",
			wantRepoCall:     true,
			wantRepoOrderNum: "79927398713",
		},
		{
			name:             "trims order number before validation and repository call",
			orderNumber:      "\n 79927398713 \t",
			wantRepoCall:     true,
			wantRepoOrderNum: "79927398713",
		},
		{
			name:        "empty order number",
			orderNumber: "",
			wantErr:     domain.ErrOrderEmpty,
		},
		{
			name:        "blank order number",
			orderNumber: "  \n\t  ",
			wantErr:     domain.ErrOrderEmpty,
		},
		{
			name:        "invalid luhn order number",
			orderNumber: "79927398710",
			wantErr:     domain.ErrInvalidOrderNumber,
		},
		{
			name:             "repository error is wrapped",
			orderNumber:      "79927398713",
			repoErr:          domain.ErrOrderAlreadyUploadedByAnotherUser,
			wantErr:          domain.ErrOrderAlreadyUploadedByAnotherUser,
			wantRepoCall:     true,
			wantRepoOrderNum: "79927398713",
		},
		{
			name:             "repository generic error is wrapped",
			orderNumber:      "79927398713",
			repoErr:          repoGenericErr,
			wantErr:          repoGenericErr,
			wantRepoCall:     true,
			wantRepoOrderNum: "79927398713",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repository := mocks.NewMockGophermartRepository(ctrl)
			service := gophermartService{repository: repository}

			if tt.wantRepoCall {
				repository.EXPECT().
					CreateOrder(gomock.Any(), userID, tt.wantRepoOrderNum).
					Return(tt.repoErr)
			}

			err := service.CreateOrder(context.Background(), userID, tt.orderNumber)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.wantErr))
				return
			}

			require.NoError(t, err)
		})
	}
}
