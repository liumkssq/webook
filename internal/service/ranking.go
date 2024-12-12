package service

import "context"

type RankingService interface {
	TopN(ctx context.Context) error
}

type BatchRankingService struct {
	artSvc  ArticleService
	intrSvc InteractiveService
}

func NewBatchRankingService(artSvc ArticleService, intrSvc InteractiveService) RankingService {
	return &BatchRankingService{artSvc: artSvc, intrSvc: intrSvc}
}

func (svc *BatchRankingService) TopN(ctx context.Context) error {
	for {
		svc.artSvc.ListPub(ctx, 0, 10)
		svc.intrSvc.GetById(ctx, 0)
	}
}
