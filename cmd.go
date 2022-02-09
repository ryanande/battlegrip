package battlegrip

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var UICmd = &cobra.Command{
  Use:   "ui",
  Short: "Launch a web based UI tool that uses the CLI",
  Long:  `Additional long description here...`,
  Run: func(cmd *cobra.Command, args []string) {
	err := Serve(cmd.Root())
	if err != nil {
		log.Fatal(err)
	}
  },
}