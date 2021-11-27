package cmd

import (
	"os"
	"os/signal"
	"syscall"

  "github.com/pkg/errors"
  "github.com/spf13/cobra"

  "github.com/moratsam/cluster/internal/broadcaster"
  "github.com/moratsam/cluster/internal/node"
)

func runCmdValidator(cmd *cobra.Command, args []string) error {
	if len(args) != 1{
		return errors.New("requires precisely 1 arg")
	}
	if args[0] == "broadcaster" || args[0] == "node" {
		return nil
	}
	return errors.New("argument must be \"broadcaster\" or \"node\"")
}

func (c *cli) runCmdRun(cmd *cobra.Command, args []string) error{
	if args[0] == "broadcaster" {
		_, err := broadcaster.NewAgent(c.cfg.broadcasterConfig)
		if err != nil {
			return err
		}
	} else{
		_, err := node.NewAgent(c.cfg.nodeConfig)
		if err != nil {
			return err
		}
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc 
	return nil //todo gracefully shutdown by closing grpc streams
}

func (c *cli) newRunCmd() *cobra.Command {

	runCmd := &cobra.Command{
		Use:		"run [broadcaster, node]",
		Short:	"run something",
		Args:		runCmdValidator,
		RunE:		c.runCmdRun,
	}

	return runCmd
}
