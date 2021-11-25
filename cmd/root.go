package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:		"cluster",
	Short:	"Interact with the cluster",
	Long:		`The only thing that matters is the silent endurance of a few, whose impassible 
presence as “stone guests” helps to create new relationships, new distances,
new values, and helps to construct a pole that, although it will certainly 
not prevent this world inhabited by the distracted and restless from being 
what it is, will still help to transmit to someone the sensation of the truth —
a sensation that could become for them the principle of a liberating crisis.`,
}

func Execute() {
	rootCmd.AddCommand(newRunCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

