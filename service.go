package main

import (
	"context"
	"log"

	"github.com/e-commerce-microservices/admin-service/pb"
	"github.com/e-commerce-microservices/admin-service/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reportService struct {
	authClient    pb.AuthServiceClient
	orderClient   pb.OrderServiceClient
	productClient pb.ProductServiceClient
	repo          *repository.Queries
	pb.UnimplementedAdminServiceServer
}

func (srv reportService) Ping(context.Context, *empty.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Message: "pong",
	}, nil
}
func (srv reportService) CreateReport(ctx context.Context, req *pb.CreateReportRequest) (*pb.CreateReportResponse, error) {
	// authen
	log.Println("check authentication")
	md, _ := metadata.FromIncomingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// check order is handled
	resp, err := srv.orderClient.CheckOrderIsHandled(ctx, &pb.CheckOrderIsHandledRequest{
		ProductId: req.GetProductId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	if !resp.GetIsBought() {
		return nil, status.Error(codes.PermissionDenied, "Sản phẩm chưa được mua")
	}

	err = srv.repo.CreateReport(ctx, repository.CreateReportParams{
		ProductID:   req.GetProductId(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &pb.CreateReportResponse{
		Message: "Báo cáo sản phẩm thành công, vui lòng chờ admin xử lý",
	}, nil
}

func (srv reportService) HandleReport(ctx context.Context, req *pb.HandleReportRequest) (*pb.HandleReportResponse, error) {
	// check admin role
	md, _ := metadata.FromIncomingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := srv.authClient.AdminAuthorization(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Không có quyền để thực hiện")
	}

	report, err := srv.repo.HandleReport(ctx, req.GetReportId())
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Xử lý không thành công")
	}

	// delete product
	resp, err := srv.productClient.DeleteProductByAdmin(ctx, &pb.DeleteProductByAdminRequest{
		ProductId: report.ProductID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Xử lý không thành công")
	}

	return &pb.HandleReportResponse{
		Message: resp.GetMessage(),
	}, nil
}

func (srv reportService) DeleteReport(ctx context.Context, req *pb.HandleReportRequest) (*pb.HandleReportResponse, error) {
	// check admin role
	md, _ := metadata.FromIncomingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := srv.authClient.AdminAuthorization(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Không có quyền để thực hiện")
	}

	err = srv.repo.DeleteReport(ctx, req.GetReportId())
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &pb.HandleReportResponse{
		Message: "Xóa report thành công",
	}, nil
}
