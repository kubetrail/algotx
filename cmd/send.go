/*
Copyright © 2022 kubetrail.io authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/kubetrail/algotx/pkg/flags"
	"github.com/kubetrail/algotx/pkg/run"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send micro algos to an address",
	Long: `This command lets you sign a transaction on
Algorand blockchain that sends a specific value of 
Algo tokens (coins) to a receiver address`,
	RunE: run.Send,
}

func init() {
	rootCmd.AddCommand(sendCmd)
	f := sendCmd.Flags()

	f.String(flags.Addr, "", "Address of the receiver")
	f.String(flags.Key, "", "Private key of the sender")
	f.Uint64(flags.Amount, 0, "Number of micro Algos to send")
	f.String(flags.Memo, "", "Transaction memo")
}
