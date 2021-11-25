package cmd

import (
  "fmt"
  "time"

  "github.com/pkg/errors"
  "github.com/spf13/cobra"
  "github.com/spf13/pflag"
  "github.com/spf13/viper"
)

func initialiseConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.AutomaticEnv() //read env vars
	bindFlags(cmd, v) //bind env vars to cobra flags
	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper){
	cmd.Flags().VisitAll(func(f *pflag.Flag){
		if v.IsSet(f.Name){
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func runCmdValidator(cmd *cobra.Command, args []string) error {
	if len(args) != 1{
		return errors.New("requires precisely 1 arg")
	}
	if args[0] == "broadcaster" || args[0] == "node" {
		return nil
	}
	return errors.New("argument must be \"broadcaster\" or \"node\"")
}

func runCmdRun(cmd *cobra.Command, args []string){
	for {
		time.Sleep(1 * time.Second)
		fmt.Println(args[0])
	}
}

func newRunCmd() *cobra.Command {
	bcaster_addr := ""

	runCmd := &cobra.Command{
		Use: "run [broadcaster, node]",
		Short: "run something",
		Args: runCmdValidator,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
				return initialiseConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args[]string){
			for {
				time.Sleep(1 * time.Second)
				fmt.Println("bcaster_addr: ", bcaster_addr)
			}
			
		},
	}

	runCmd.Flags().StringVarP(&bcaster_addr, "BCASTER_ADDR", "b", "", "address of the broadcaster")
	return runCmd
}
