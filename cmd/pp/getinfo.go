package main

import (
	"fmt"

	"github.com/mniak/bcpinpad/highlevel"
	"github.com/spf13/cobra"
)

var (
	getInfo = &cobra.Command{
		Use: "getinfo",
		Aliases: []string{
			"get-info",
			"getInfo",
			"GetInfo",
		},
		Short: "Get info about the pinpad",
		RunE: func(cmd *cobra.Command, args []string) error {
			pinpad, err := NewPinpad(cmd, args)
			if err != nil {
				return err
			}
			params := highlevel.GetInfoParams{
				Type: 0,
			}
			result, err := pinpad.GetInfo(params)
			if err != nil {
				return err
			}
			if params.Type == 0 {
				fmt.Printf("Manufacturer: %s\n", result.Name)
			} else {
				fmt.Printf("Acquirer: %s\n", result.Name)
			}
			return nil
		},
	}
)
