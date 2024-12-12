package grpc

import "context"

type Server struct {
	UnimplementedUserServiceServer
}

func (s Server) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	return &GetByIdResp{
		User: &User{
			Id:   1,
			Name: "test",
		},
	}, nil
}
