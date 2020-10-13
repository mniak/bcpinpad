package main

import (
	"github.com/mniak/bcpinpad/highlevel"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "pp",
	}
)

func init() {
	rootCmd.AddCommand(getInfo)
	rootCmd.AddCommand(display)
}

func main() {
	rootCmd.Execute()
}

func NewPinpad(cmd *cobra.Command, args []string) (highlevel.Pinpad, error) {
	pp, err := highlevel.OpenSerial("COM4")
	return pp, err
}
