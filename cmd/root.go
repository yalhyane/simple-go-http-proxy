/*
Copyright © 2023 Yacine Alhyane y.alhyane@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yalhyane/simple-http-proxy/internal"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simple-http-proxy",
	Short: "A simple http proxy",
	Long: `This is a simple HTTP proxy project built in Golang.
This was created for learning purposes to understand how HTTP proxies work.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	RunE: runE,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var proxyConfig internal.CustomProxyConfig

func init() {
	proxyConfig = internal.CustomProxyConfig{}
	rootCmd.Flags().StringVarP(&proxyConfig.Addr, "addr", "a", "0.0.0.0:8889", "Proxy listen address")
	rootCmd.Flags().DurationVarP(&proxyConfig.TargetTimeout, "target-timeout", "T", 10*time.Second, "Target timeout")
	rootCmd.Flags().BoolVar(&proxyConfig.Verbose, "verbose", true, "Proxy verbosity")
}

func runE(cmd *cobra.Command, args []string) error {
	if err := proxyConfig.ValidateE(); err != nil {
		return err
	}
	proxyConfig.FillDefaults()

	proxy := internal.CustomProxy{
		Config: proxyConfig,
	}

	proxy.Start()

	return nil
}
