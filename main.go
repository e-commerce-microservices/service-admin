package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/e-commerce-microservices/admin-service/pb"
	"github.com/e-commerce-microservices/admin-service/repository"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	// postgres driver
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// init user db connection
	pgDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWD"), os.Getenv("DB_DBNAME"),
	)
	reportDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer reportDB.Close()
	if err := reportDB.Ping(); err != nil {
		log.Fatal("can't ping to user db", err)
	}

	grpcServer := grpc.NewServer()

	authConn, err := grpc.Dial("auth-service:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't dial auth service", err)
	}
	authClient := pb.NewAuthServiceClient(authConn)

	productConn, err := grpc.Dial("product-service:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't dial product service", err)
	}
	productClient := pb.NewProductServiceClient(productConn)

	orderConn, err := grpc.Dial("order-service:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't dial product service", err)
	}
	orderClient := pb.NewOrderServiceClient(orderConn)

	queries := repository.New(reportDB)
	reportService := reportService{
		authClient:    authClient,
		orderClient:   orderClient,
		productClient: productClient,
		repo:          queries,
	}
	pb.RegisterAdminServiceServer(grpcServer, reportService)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}
	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}
}
