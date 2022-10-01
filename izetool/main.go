package main

import "com.samderlust/izetool/izetool/command"

func main() {

	rootCmd := command.RootCommand()

	// rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
	// 	return nil
	// }

	rootCmd.Execute()
}
