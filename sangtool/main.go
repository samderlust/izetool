package main

import "log"

func main() {
	log.SetFlags(0)

	rootCmd := rootCommand()

	// rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
	// 	return nil
	// }

	rootCmd.Execute()
}
