package node

import (
	"context"

	"google.golang.org/grpc"
	"go.uber.org/zap"

	api_msg "github.com/moratsam/cluster/api/v1/msg"
	api_node "github.com/moratsam/cluster/api/v1/node"
	api_queue "github.com/moratsam/cluster/api/v1/queue"
)

var _ api_node.NodeServer = (*grpcServer)(nil)

type grpcServer struct {
	api_node.UnimplementedNodeServer

	logger			*zap.Logger
	bcast_client	api_queue.QueueClient

}

func NewGRPCServer(broadcaster_addr string) (*grpc.Server, error){
	gsrv := grpc.NewServer()
	srv, err := newgrpcServer(broadcaster_addr)
	if err != nil {
		return nil, err
	}
	api_node.RegisterNodeServer(gsrv, srv)
	return gsrv, nil
}

func newgrpcServer(broadcaster_addr string) (*grpcServer, error) {

	clientOptions := []grpc.DialOption{grpc.WithInsecure()}
	cc, err := grpc.Dial(broadcaster_addr, clientOptions...)
	if err != nil {
		return nil, err
	}
	bcast_client := api_queue.NewQueueClient(cc)

	srv := &grpcServer{
		logger:			zap.L().Named("node_server"),
		bcast_client:	bcast_client,
	}
	return srv, nil
}

//forward message to broadcaster
func (s *grpcServer) broadcast(msg *api_msg.Msg) (*api_queue.Ack, error){
	ctx := context.Background()

	ack, err := s.bcast_client.Publish(
		ctx, 
		msg,
	)
	if err != nil {
		s.logger.Error("failed to forward message to broadcaster", zap.Error(err))
		return nil, err
	}
	if ack.Ok != true {
		s.logger.Error("broadcaster Ack is not true")
	}
	return ack, nil
}

//someone published something, so republish it to subscribers and send ack
func (s *grpcServer) Publish(ctx context.Context, msg *api_msg.Msg) (*api_node.Ack, error) {

		ack, err := s.broadcast(msg)
		if err != nil {
			return nil, err
		}

		res := &api_node.Ack{Ok: ack.Ok}
		return res, nil
}

//someone sent a subscription request, so subscribe to broadcaster and forward back messages
func (s *grpcServer) Subscribe(
	req		*api_node.SubscriptionRequest,
	stream	api_node.Node_SubscribeServer,
) error {

	ctx := context.Background()
	msg_type := req.GetType()

	//subscribe to broadcaster
	bcaster_stream, err := s.bcast_client.Subscribe(
		ctx,
		&api_queue.SubscriptionRequest{
			Type: msg_type,
		},
	)
	if err != nil {
		s.logger.Error("error subscribing to broadcaster", zap.Error(err))
		return err
	}

	for {
		select{
		case <- bcaster_stream.Context().Done():
			return nil
		default:
			msg, err := bcaster_stream.Recv()
			if err != nil {
				s.logger.Error("error while receiving from broadcaster stream", zap.Error(err))
			}
			if err := stream.Send(msg); err != nil {
				s.logger.Error("failed to forward message back to client", zap.Error(err))
			}
		}
	}
}
