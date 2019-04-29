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
	"github.com/posener/complete"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "balance|chains|transactions",
		Long: `
Get balance, transaction, or issuance data about an existing FAT Chain.

The fatd API is used to lookup information about FAT chains. Thus fat-cli can
only return data about chains that the instance of fatd is tracking. The fatd
API must be trusted to ensure the security and validity of returned data.
`[1:],
	}
	rootCmd.AddCommand(cmd)
	rootCmplCmd.Sub["get"] = getCmplCmd
	rootCmplCmd.Sub["help"].Sub["get"] = complete.Command{Sub: complete.Commands{}}
	generateCmplFlags(cmd, getCmplCmd.Flags)
	return cmd
}()

var getCmplCmd = complete.Command{
	Flags: mergeFlags(apiCmplFlags, tokenCmplFlags),
	Sub:   complete.Commands{},
}