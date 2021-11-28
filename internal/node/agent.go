package node

import (
	"net"

	_"google.golang.org/grpc"
	"go.uber.org/zap"

)

type Config struct {
	BindAddr					string
	BroadcasterAddr		string
}

type Agent struct {
	Config
}

func NewAgent(config Config) (*Agent, error) {
	a := &Agent{ Config: config }
	setup := []func() error{
		a.setupLogger,
		a.setupServer,
	}
	for _, fn := range setup {
		if err := fn(); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (a *Agent) setupLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	logger.Debug("node agent",
		zap.String("BindAddr", a.Config.BindAddr),
		zap.String("BroadcasterAddr", a.Config.BroadcasterAddr),
	)
	return nil
}

func (a *Agent) setupServer() error {
	addr, err := net.ResolveTCPAddr("tcp", a.Config.BindAddr)
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}
	addr, err = net.ResolveTCPAddr("tcp", a.Config.BroadcasterAddr)
	if err != nil {
		return err
	}
	server, err := NewGRPCServer(addr.String())
	if err != nil {
		return err
	}
	
	go func(){
		server.Serve(l)
	}()

	return nil
}
