package run

import (
	"github.com/kubetrail/algotx/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type persistentFlagValues struct {
	OutputFormat string `json:"outputFormat,omitempty"`
	RPCEndpoint  string `json:"rpcEndPoint,omitempty"`
	RPCToken     string `json:"rpcToken,omitempty"`
}

func getPersistentFlags(cmd *cobra.Command) persistentFlagValues {
	rootCmd := cmd.Root().PersistentFlags()

	_ = viper.BindPFlag(flags.OutputFormat, rootCmd.Lookup(flags.OutputFormat))
	_ = viper.BindPFlag(flags.RPCEndpoint, rootCmd.Lookup(flags.RPCEndpoint))
	_ = viper.BindPFlag(flags.RPCToken, rootCmd.Lookup(flags.RPCToken))

	outputFormat := viper.GetString(flags.OutputFormat)
	rpcEndpoint := viper.GetString(flags.RPCEndpoint)
	rpcToken := viper.GetString(flags.RPCToken)

	return persistentFlagValues{
		OutputFormat: outputFormat,
		RPCEndpoint:  rpcEndpoint,
		RPCToken:     rpcToken,
	}
}
