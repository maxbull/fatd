// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/Factom-Asset-Tokens/fatd/factom"
	"github.com/Factom-Asset-Tokens/fatd/srv"
	"github.com/posener/complete"
	"github.com/spf13/cobra"
)

var addresses []factom.FAAddress

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:                   "balance ADDRESS...",
	Aliases:               []string{"balances"},
	DisableFlagsInUseLine: true,
	Short:                 "Get the balances for addresses",
	Long: `Get the balances of the listed addresses.

Queries fatd for the balances for each ADDRESS for the specified FAT Chain.

Required flags: --chainid or --tokenid and --identity`,
	Args:    getBalanceArgs,
	PreRunE: validateChainID,
	Run:     getBalance,
}

var balanceCmplCmd = complete.Command{
	Flags: rootCmplCmd.Flags,
	Args:  PredictFAAddresses,
}

func init() {
	getCmd.AddCommand(balanceCmd)
	getCmplCmd.Sub["balance"] = balanceCmplCmd
}

func getBalanceArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	addresses = make([]factom.FAAddress, len(args))
	duplicate := make(map[factom.FAAddress]struct{}, len(args))
	for i := range addresses {
		adr := &addresses[i]
		if err := adr.Set(args[i]); err != nil {
			return err
		}
		if _, ok := duplicate[*adr]; ok {
			return fmt.Errorf("duplicate: %v", adr)
		}
		duplicate[*adr] = struct{}{}
	}
	return nil
}

func getBalance(cmd *cobra.Command, _ []string) {
	var params srv.ParamsGetBalance
	params.ChainID = &ChainID

	balances := make([]uint64, len(addresses))
	for i, adr := range addresses {
		params.Address = &adr
		if err := FATClient.Request("get-balance", params, &balances[i]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	for i, adr := range addresses {
		fmt.Println(adr, balances[i])
	}
}
