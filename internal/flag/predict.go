// MIT License
//
// Copyright 2018 Canonical Ledgers, LLC
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package flag

import (
	"context"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/posener/complete"
)

func predictAddress(fa bool, num int, flagName, suffix string) complete.PredictFunc {
	if len(flagName) == 0 {
		return func(a complete.Args) []string {
			// Count the number of complete arguments that are not flags.
			argc := len(a.Completed[1:])
			for _, arg := range a.Completed[1:] {
				if string(arg[0]) == "-" {
					argc--
				}
			}
			if len(suffix) > 0 && len(a.Last) > 0 &&
				a.Last[len(a.Last)-1:len(a.Last)] == suffix {
				return nil
			}
			if argc < num {
				adrs := listAddresses(fa)
				if len(suffix) > 0 {
					for i := range adrs {
						adrs[i] += suffix
					}
				}
				return adrs
			}
			return nil
		}
	}
	return func(a complete.Args) []string {
		// Count the number of complete arguments that are not flags.
		argc := 0
		for i := len(a.Completed) - 1; i > 0; i-- {
			arg := a.Completed[i]
			if string(arg) == flagName {
				break
			}
			argc++
		}
		if len(suffix) > 0 && len(a.Last) > 0 &&
			a.Last[len(a.Last)-1:len(a.Last)] == suffix {
			return nil
		}
		if argc < num {
			adrs := listAddresses(fa)
			if len(suffix) > 0 {
				for i := range adrs {
					adrs[i] += suffix
				}
			}
			return adrs
		}
		return nil
	}
}

func listAddresses(fa bool) []string {
	parseWalletFlags()
	fss, ess, err := FactomClient.GetPrivateAddresses(context.Background())
	if err != nil {
		os.Exit(6)
	}
	var adrStrs []string
	if fa {
		adrStrs = make([]string, len(fss))
		for i, fs := range fss {
			adrStrs[i] = fs.FAAddress().String()
		}
	} else {
		adrStrs = make([]string, len(ess))
		for i, es := range ess {
			adrStrs[i] = es.ECAddress().String()
		}
	}
	return adrStrs
}

var cliFlags *flag.FlagSet

// Parse any previously specified factom-cli options required for connecting to
// factom-walletd
func parseWalletFlags() {
	if cliFlags != nil {
		// We already parsed the flags.
		return
	}
	// Using flag.FlagSet allows us to parse a custom array of flags
	// instead of this programs args.
	cliFlags = flag.NewFlagSet("", flag.ContinueOnError)
	cliFlags.StringVar(&FactomClient.WalletdServer, "w", "localhost:8089", "")
	cliFlags.StringVar(&FactomClient.Walletd.User, "walletuser", "", "")
	cliFlags.StringVar(&FactomClient.Walletd.Password, "walletpassword", "", "")

	// flags.Parse will print warnings if it comes across an unrecognized
	// flag. We don't want this so we temprorarily redirect everything to
	// /dev/null before we call flags.Parse().
	stdout, stderr := os.Stdout, os.Stderr
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		os.Exit(5)
	}
	os.Stdout, os.Stderr = devNull, devNull

	// The current command line being typed is stored in the environment
	// variable COMP_LINE. We split on spaces and discard the first in the
	// list because it is the program name `factom-cli`.
	cliFlags.Parse(strings.Fields(os.Getenv("COMP_LINE"))[1:])

	// Restore stdout and stderr.
	os.Stdout, os.Stderr = stdout, stderr

	// We need a short timeout or the CLI completion will hang.
	FactomClient.Walletd.Timeout = time.Second / 2
}