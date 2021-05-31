package battlegrip

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
)


func TestCommandEndpoint(t *testing.T)  {

	rootCmd := &cobra.Command{
		Use:   "pcs",
		Short: "pcs cli tool.",
		Long: `Wheel Snipe Celey!`,
	}

	rootCmd.Execute()

	var uiCmd = &cobra.Command{
		Use:   "ui",
		Short: "Launch a web based UI tool that uses the CLI",
		Long:  `Additional long description here...`,
		Run: func(cmd *cobra.Command, args []string) {
		  
	  
		  err := Serve(cmd)
		  if err != nil {
			  fmt.Errorf("%v", err)
		  }
		},
	  }

	  uiCmd.ExecuteC()
}