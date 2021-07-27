package main

func main() {

	rootCmd := rootCommand()

	// rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
	// 	return nil
	// }

	rootCmd.Execute()
}
