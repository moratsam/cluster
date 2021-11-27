package tests

import(
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
	"google.golang.org/grpc"

	api_msg "github.com/moratsam/cluster/api/v1/msg"
	api_node "github.com/moratsam/cluster/api/v1/node"
	"github.com/moratsam/cluster/internal/broadcaster"
	"github.com/moratsam/cluster/internal/node"
)

func get_client(t *testing.T, server_addr string) api_node.NodeClient {
	clientOptions := []grpc.DialOption{grpc.WithInsecure()}
	cc, err := grpc.Dial(server_addr, clientOptions...)
	require.NoError(t, err)
	client := api_node.NewNodeClient(cc)
	return client
}

func TestAgents(t *testing.T) {
	ports := dynaport.Get(3)
	broadcaster_addr	:= fmt.Sprintf("%s:%d", "127.0.0.1", ports[0])
	node1_addr			:= fmt.Sprintf("%s:%d", "127.0.0.1", ports[1])
	node2_addr			:= fmt.Sprintf("%s:%d", "127.0.0.1", ports[2])

	_, err := broadcaster.NewAgent(
		broadcaster.Config{
			BindAddr: broadcaster_addr,
		},
	)
	require.NoError(t, err)

	var node_agents []*node.Agent
	for i:=0; i<2; i++{
		agent, err := node.NewAgent(
			node.Config{
				BindAddr:			fmt.Sprintf("%s:%d", "127.0.0.1", ports[i+1]),
				BroadcasterAddr:	broadcaster_addr,
			},
		)
		require.NoError(t, err)

		node_agents = append(node_agents, agent)
	}

	client1 := get_client(t, node1_addr)
	client2 := get_client(t, node2_addr)

	ctx := context.Background()

	//client1 gets sub stream
	sub_stream1, err := client1.Subscribe(
		ctx,
		&api_node.SubscriptionRequest{
			Type: api_msg.MsgType_VANILLA,
		},
	)
	require.NoError(t, err)

	//client2 gets sub stream
	sub_stream2, err := client2.Subscribe(
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
	msg, err := sub_stream1.Recv()
	require.NoError(t, err)
	if msg.Type != api_msg.MsgType_VANILLA {
		t.Fatalf("got msg type: %v, expected: %v", msg.Type, api_msg.MsgType_VANILLA)
	}
	if msg.Data != "kurbarija" {
		t.Fatalf("got msg data: %v, expected: %v", msg.Data, "kurbarija")
	}

	//client2 receives own publication
	msg, err = sub_stream2.Recv()
	require.NoError(t, err)
	if msg.Type != api_msg.MsgType_VANILLA {
		t.Fatalf("got msg type: %v, expected: %v", msg.Type, api_msg.MsgType_VANILLA)
	}
	if msg.Data != "kurbarija" {
		t.Fatalf("got msg data: %v, expected: %v", msg.Data, "kurbarija")
	}
}

