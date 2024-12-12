package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestBatchRankingService_TopN(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (ArticleService, InteractiveService)
		wantErr error
	}{
		{
			name: "normal",
			mock: func(ctrl *gomock.Controller) (ArticleService, InteractiveService) {

			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			artSvc, intrSvc := tc.mock(ctrl)
			svc := NewBatchRankingService(artSvc, intrSvc)
			err := svc.TopN(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
