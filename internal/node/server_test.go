package node

import (
	"context"
	"net"
	"testing"
	_"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	broadcaster "github.com/moratsam/cluster/internal/broadcaster"
	api_msg "github.com/moratsam/cluster/api/v1/msg"
	api_node "github.com/moratsam/cluster/api/v1/node"
)

func TestServer(t *testing.T){
	for scenario, fn := range map[string]func(
		t *testing.T,
		client1 api_node.NodeClient,
		client2 api_node.NodeClient,
	){
		"subscribe and hear broadcast of own publication": testSubscribeToSelf,
		"subscribe and hear broadcast of another's publication": testSubscribeAndHearAnother,
	} {
		t.Run(scenario, func(t *testing.T){
			client1, client2, teardown := setupTest(t, nil)
			defer teardown()
			fn(t, client1, client2)
		})
	}
}

func setupTest(t *testing.T, fn func()) (
	client	api_node.NodeClient,
	client2	api_node.NodeClient,
	teardown	func(),
) {
	t.Helper()

	//fire up broadcaster
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	broadcaster_addr := l.Addr().String()
	broadcaster, err := broadcaster.NewGRPCServer()
	require.NoError(t, err)
	go func(){
		broadcaster.Serve(l)
	}()

	//fire up two node servers
	l1, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	clientOptions := []grpc.DialOption{grpc.WithInsecure()}
	cc1, err := grpc.Dial(l1.Addr().String(), clientOptions...)
	require.NoError(t, err)

	server1, err := NewGRPCServer(broadcaster_addr)
	require.NoError(t, err)

	go func(){
		server1.Serve(l1)
	}()

	client1 := api_node.NewNodeClient(cc1)


	l2, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	cc2, err := grpc.Dial(l2.Addr().String(), clientOptions...)
	require.NoError(t, err)

	server2, err := NewGRPCServer(broadcaster_addr)
	require.NoError(t, err)

	go func(){
		server2.Serve(l2)
	}()

	client2 = api_node.NewNodeClient(cc2)
	
	return client1, client2, func(){
		server1.Stop()
		server2.Stop()
		broadcaster.Stop()
		cc1.Close()
		cc2.Close()
		l.Close()
		l1.Close()
		l2.Close()
	}
}

func testSubscribeToSelf(t *testing.T, client, _ api_node.NodeClient){
	ctx := context.Background()
	//get sub stream
	sub_stream, err := client.Subscribe(
		ctx,
		&api_node.SubscriptionRequest{
			Type: api_msg.MsgType_VANILLA,
		},
	)
	require.NoError(t, err)

	//publish msg, receive ack
	ack, err := client.Publish(
		ctx, 
		&api_msg.Msg{
			Type: api_msg.MsgType_VANILLA,
			Data: "AI reconquista",
		},
	)
	require.NoError(t, err)
	if ack.Ok != true {
		t.Fatalf("got ok: %v, expected: %v", ack.Ok, true)
	}

	//receive msg from sub_stream
	msg, err := sub_stream.Recv()
	require.NoError(t, err)
	if msg.Type != api_msg.MsgType_VANILLA {
		t.Fatalf("got msg type: %v, expected: %v", msg.Type, api_msg.MsgType_VANILLA)
	}
	if msg.Data != "AI reconquista" {
		t.Fatalf("got msg data: %v, expected: %v", msg.Data, "AI reconquista")
	}
}

func testSubscribeAndHearAnother(t *testing.T, client1, client2 api_node.NodeClient){
	ctx := context.Background()

	//client1 gets sub stream
	sub_stream, err := client1.Subscribe(
		ctx,
		&api_node.SubscriptionRequest{
			Type: api_msg.MsgType_VANILLA,
		},
	)
	require.NoError(t, err)

	//client2 published msg
	ack, err := client2.Publish(
		ctx,
		&api_msg.Msg{
			Type: api_msg.MsgType_VANILLA,
			Data: "kurbarija",
		},
	)
	require.NoError(t, err)
	if ack.Ok != true {
		t.Fatalf("got ok %v, expected %v", ack.Ok, true)
	}

	//client1 receives client2's publication
	msg, err := sub_stream.Recv()
	require.NoError(t, err)
	if msg.Type != api_msg.MsgType_VANILLA {
		t.Fatalf("got msg type: %v, expected: %v", msg.Type, api_msg.MsgType_VANILLA)
	}
	if msg.Data != "kurbarija" {
		t.Fatalf("got msg data: %v, expected: %v", msg.Data, "kurbarija")
	}
}

