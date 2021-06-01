package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:           "cmstore",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	rootFlags = struct {
		dir       string
		namespace string
		name      string
	}{}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootFlags.dir, "dir", "d", "", "")
	rootCmd.MarkPersistentFlagRequired("dir")
	rootCmd.PersistentFlags().StringVar(&rootFlags.namespace, "namespace", "", "")
	rootCmd.MarkPersistentFlagRequired("namespace")
	rootCmd.PersistentFlags().StringVar(&rootFlags.name, "name", "", "")
	rootCmd.MarkPersistentFlagRequired("name")
}
