package run

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/kubetrail/algotx/pkg/flags"
	"github.com/kubetrail/bip32/pkg/keys"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Balance(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.Addr, cmd.Flag(flags.Addr))
	addr := viper.GetString(flags.Addr)

	prompt, err := prompts.Status()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	if len(addr) == 0 {
		if len(args) == 0 {
			if prompt {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter address: "); err != nil {
					return fmt.Errorf("failed to write to output: %w", err)
				}
			}
			addr, err = keys.Read(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read pub addr from input: %w", err)
			}
		} else {
			addr = args[0]
		}
	}

	algodClient, err := algod.MakeClient(persistentFlags.RPCEndpoint, persistentFlags.RPCToken)
	if err != nil {
		return fmt.Errorf("failed to make algod client: %w", err)
	}

	accountInfo, err := algodClient.AccountInformation(addr).Do(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatNative:
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), accountInfo.Amount); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	case flags.OutputFormatYaml:
		b, err := yaml.Marshal(accountInfo)
		if err != nil {
			return fmt.Errorf("failed to serialize output to yaml: %w", err)
		}

		if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
			return fmt.Errorf("failed to write yaml to output: %w", err)
		}
	case flags.OutputFormatJson:
		b, err := json.Marshal(accountInfo)
		if err != nil {
			return fmt.Errorf("failed to serialize output to json: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(b)); err != nil {
			return fmt.Errorf("failed to write json to output: %w", err)
		}
	default:
		return fmt.Errorf("invalid output format")
	}

	return nil
}
