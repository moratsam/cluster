package broadcaster

import (
	"net"

	_"google.golang.org/grpc"
	"go.uber.org/zap"

)

type Config struct {
	BindAddr		string
}

type Agent struct {
	logger *zap.Logger
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
	logger.Debug("broadcaster agent", zap.String("BindAddr", a.Config.BindAddr))
	a.logger = logger
	return nil
}

func (a *Agent) setupServer() error {
	addr, err := net.ResolveTCPAddr("tcp", a.Config.BindAddr) //resolve hostname to IP
	if err != nil {
		return err
	}
	a.logger.Debug("", zap.String("addr", addr.String()))
	l, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}
	server, err := NewGRPCServer()
	if err != nil {
		return err
	}
	
	go func(){
		server.Serve(l)
	}()

	return nil
}
