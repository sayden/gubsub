package grpcservice

import (
	"fmt"
	"net"
	"time"

	"errors"
	log "github.com/Sirupsen/logrus"
	serfclient "github.com/hashicorp/serf/client"
	"github.com/sayden/gubsub/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type rpc_server struct {
	Port       int
	Dispatcher types.Dispatcher
}

var RPC rpc_server

func NewGRPCServer(d types.Dispatcher) {
	RPC = rpc_server{
		Port:       5123,
		Dispatcher: d,
	}

	log.Info("RPC Server created")

	go func() {
		time.Sleep(2 * time.Second)
		RPC.StartServer()
	}()
}

//NewMessage is the implementation to receive a new message across the cluster
func (s *rpc_server) NewMessage(ctx context.Context, in *GubsubMessage) (*GubsubReply, error) {
	log.Info("gRPC message received for topic", in.T)
	m := types.Message{
		Data:      &in.M,
		Topic:     &in.T,
		Timestamp: time.Now(),
	}

	if s.Dispatcher == nil {
		log.Error("Dispatcher is nil")
		return nil, errors.New("dispatcher is nil")
	}
	s.Dispatcher.DispatchMessageLocal(&m)
	return &GubsubReply{StatusCode: 0}, nil
}

func (s *rpc_server) StartServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Error("Failed starting server: failed to listen:", err)
	}

	server := grpc.NewServer()
	RegisterMessageServiceServer(server, &rpc_server{})
	server.Serve(lis)
}

func (s *rpc_server) SendMessageInCluster(m *types.Message, servers []serfclient.Member) error {
	codes := make([]int32, len(servers))
	fmt.Printf(" %d %s \n", len(servers), string(*m.Data))
	for _, server := range servers {
		fmt.Printf(" %s \n", server.Addr)
		statusCode, err := s.SendMessage(m, server)
		if err != nil {
			log.Error("Error trying to send message to member %s", server.Addr.String())
		}

		codes = append(codes, statusCode)
	}

	for _, c := range codes {
		if c != 0 {
			return errors.New("Not all messages have been delivered correctly")
		}
	}

	return nil
}

func (s *rpc_server) SendMessage(m *types.Message, server serfclient.Member) (int32, error) {
	log.Info("Sending a message across cluster")

	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", server.Addr.String(), s.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewMessageServiceClient(conn)

	r, err := c.NewMessage(context.Background(), &GubsubMessage{
		M: *m.Data,
		T: *m.Topic,
	})

	if err != nil {
		return -1, err
	}

	return r.StatusCode, nil
}
