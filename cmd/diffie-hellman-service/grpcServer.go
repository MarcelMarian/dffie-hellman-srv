package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type dhGrpcServiceServer struct {
	UnimplementedDhGrpcServiceServer
}

var (
	tls             = true
	certFile        = "/config/server_cert.pem"
	keyFile         = "/config/server_key.pem"
	port     uint32 = 10000
)

var sha string

func (s *dhGrpcServiceServer) Ping(ctx context.Context, req *PingPongMessage) (*PingPongMessage, error) {
	resp := PingPongMessage{}
	resp.Msg = "PONG"
	return &resp, nil
}

func (s *dhGrpcServiceServer) CliXchgKey(ctx context.Context, req *PKeyMessage) (*PKeyMessage, error) {
	start := time.Now()

	resp := PKeyMessage{}

	pub := new(big.Int)
	// priv := new(big.Int)

	pub, _, sha = CalcEncryptionKey(req.Pubkey)

	resp.Seqno = req.Seqno
	resp.Pubkey = pub.Bytes()

	// fmt.Println("bob   computes:", priv)
	// fmt.Println("bob   sha256:", sha)
	fmt.Printf("%d.Key xchg duration t:%v\n\n", resp.Seqno, time.Since(start).Seconds())

	return &resp, nil
}

func (s *dhGrpcServiceServer) XchgMessage(ctx context.Context, req *EncMessage) (*EncMessageResp, error) {
	start := time.Now()

	resp := EncMessageResp{}

	resp.Seqno = req.Seqno
	resp.RetCode = 0

	// fmt.Println("Data:", req.Seqno, req.Data)
	msg := decrypt(req.Data, sha)

	fmt.Printf("%d.Dec Msg: %s\n", resp.Seqno, msg)
	fmt.Printf("%d.Msg decrypt duration t:%f\n\n", resp.Seqno, time.Since(start).Seconds())

	return &resp, nil
}

// Entrypoint of GRPC Server start
func startGRPCSever() error {
	server := dhGrpcServiceServer{}

	tls = configData.GrpcServerConfig.TlsEnable
	port = configData.GrpcServerConfig.Port

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return err
	}

	log.Printf("gRPC server is listening on: 0.0.0.0:%d", port)

	var opts []grpc.ServerOption
	if tls {
		if configDataPtr.GrpcServerConfig.CertFile != "" {
			certFile = configDataPtr.GrpcServerConfig.CertFile
		}
		if configDataPtr.GrpcServerConfig.KeyFile != "" {
			keyFile = configDataPtr.GrpcServerConfig.KeyFile
		}
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Printf("Failed to generate credentials %v", err)
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)

	reflection.Register(grpcServer)

	RegisterDhGrpcServiceServer(grpcServer, &server)

	if err := grpcServer.Serve(listen); err != nil {
		log.Printf("failed to serve: %v", err)
		return err
	}

	return nil
}
