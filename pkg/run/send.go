package run

import (
	"bufio"
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/kubetrail/algotx/pkg/flags"
	"github.com/kubetrail/bip32/pkg/keys"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Send(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.Addr, cmd.Flag(flags.Addr))
	_ = viper.BindPFlag(flags.Key, cmd.Flag(flags.Key))
	_ = viper.BindPFlag(flags.Amount, cmd.Flag(flags.Amount))
	_ = viper.BindPFlag(flags.Memo, cmd.Flag(flags.Memo))

	addr := viper.GetString(flags.Addr)
	key := viper.GetString(flags.Key)
	amount := viper.GetUint64(flags.Amount)
	memo := viper.GetString(flags.Memo)

	prompt, err := prompts.Status()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	if len(addr) == 0 {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter address receiving funds: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
		addr, err = keys.Read(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("failed to read pub addr from input: %w", err)
		}
	}

	if len(key) == 0 {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter signing private key: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
		addr, err = keys.Read(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("failed to read signing private key from input: %w", err)
		}
	}

	if len(memo) == 0 && len(args) > 0 {
		memo = strings.Join(args, " ")
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return fmt.Errorf("failed to decode input key as hex string: %w", err)
	}

	if len(keyBytes) != ed25519.PrivateKeySize {
		return fmt.Errorf("invalid key size, expected %d, got %d", ed25519.PrivateKeySize, len(keyBytes))
	}

	publicKey, privateKey, err := ed25519.GenerateKey(
		bufio.NewReader(
			bytes.NewReader(
				keyBytes[:ed25519.SeedSize],
			),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to generate key pair from input key: %w", err)
	}

	account := crypto.Account{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Address:    types.Address{},
	}

	copy(account.Address[:], account.PublicKey)

	addrSender := account.Address.String()
	addrReceiver := addr

	algodClient, err := algod.MakeClient(persistentFlags.RPCEndpoint, persistentFlags.RPCToken)
	if err != nil {
		return fmt.Errorf("failed to make algod client: %w", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get suggested params: %w", err)
	}

	var minFee uint64 = transaction.MinTxnFee
	note := []byte(memo)
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := uint64(txParams.FirstRoundValid)
	lastValidRound := uint64(txParams.LastRoundValid)
	txn, err := transaction.MakePaymentTxnWithFlatFee(addrSender, addrReceiver, minFee, amount, firstValidRound, lastValidRound, note, "", genID, genHash)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	txID, signedTxn, err := crypto.SignTransaction(account.PrivateKey, txn)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	sendResponse, err := algodClient.SendRawTransaction(signedTxn).Do(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to send raw transaction: %w", err)
	}

	confirmedTxn, err := future.WaitForConfirmation(algodClient, txID, 4, cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to wait for confirmation: %w", err)
	}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatNative:
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), sendResponse); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	case flags.OutputFormatYaml:
		b, err := yaml.Marshal(confirmedTxn)
		if err != nil {
			return fmt.Errorf("failed to serialize output to yaml: %w", err)
		}

		if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
			return fmt.Errorf("failed to write yaml to output: %w", err)
		}
	case flags.OutputFormatJson:
		b, err := json.Marshal(confirmedTxn)
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
