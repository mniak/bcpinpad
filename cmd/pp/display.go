package main

import (
	"github.com/mniak/bcpinpad/highlevel"
	"github.com/spf13/cobra"
)

var (
	display = &cobra.Command{
		Use:   "display",
		Short: "Get info about the pinpad",
		RunE: func(cmd *cobra.Command, args []string) error {
			pinpad, err := NewPinpad(cmd, args)
			if err != nil {
				return err
			}
			params := highlevel.DisplayParams{
				Message: "1234567890123456" + "1234567890123456",
			}
			_, err = pinpad.Display(params)
			if err != nil {
				return err
			}
			// if params.Type == 0 {
			// 	fmt.Printf("Manufacturer: %s\n", result.Name)
			// } else {
			// 	fmt.Printf("Acquirer: %s\n", result.Name)
			// }
			return nil
		},
	}
)
