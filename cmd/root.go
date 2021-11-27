package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/moratsam/cluster/internal/broadcaster"
	"github.com/moratsam/cluster/internal/node"
)

type cfg struct {
	broadcasterConfig	broadcaster.Config
	nodeConfig			node.Config
}

type cli struct {
	cfg
}

func setupFlags(cmd *cobra.Command) error {
	cmd.Flags().String("config-file", "", "path to config file")
	cmd.Flags().String("node-bind-addr", "this should get overwritten", "addr of node grpc")
	return viper.BindPFlags(cmd.Flags())
}

func (c *cli) setupConfig(cmd *cobra.Command, args []string) error {
	config_file, err := cmd.Flags().GetString("config-file")
	if err != nil {
		log.Fatal("error getting config file flag")
		return err
	}
	viper.SetConfigFile(config_file)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("config file is necessary but was not found")
		return err
	}
	c.cfg.broadcasterConfig.BindAddr		= viper.GetString("broadcaster-bind-addr")

	c.cfg.nodeConfig.BindAddr				= viper.GetString("node-bind-addr")
	c.cfg.nodeConfig.BroadcasterAddr		= viper.GetString("broadcaster-bind-addr")
	return nil
}


func Execute() {
	cli := &cli{}

	rootCmd := &cobra.Command{
		Use:		"cluster",
		Short:	"Interact with the cluster",
		Long:		`The only thing that matters is the silent endurance of a few, whose impassible 
	presence as “stone guests” helps to create new relationships, new distances,
	new values, and helps to construct a pole that, although it will certainly 
	not prevent this world inhabited by the distracted and restless from being 
	what it is, will still help to transmit to someone the sensation of the truth —
	a sensation that could become for them the principle of a liberating crisis.`,
		PreRunE:	cli.setupConfig,
	}
	rootCmd.AddCommand(cli.newRunCmd())
	if err := setupFlags(rootCmd); err != nil {
		log.Fatal(err)
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("error executing rootCmd", err)
	}
}

