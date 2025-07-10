/*
Copyright © 2024 Dave Rawks <dave@rawks.io>

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

//go:generate stringer -type=PT ../pkg/runtime/protocol.go

package cmd

import (
	"log"
	_ "net/http/pprof"
	"os"

	logs "github.com/appscode/go/log/golog"
	"github.com/appscode/go/runtime"
	gearmand "github.com/drawks/gearhulk/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var cfg gearmand.Config
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Gearman server",
	Long: `Start the Gearman server with the specified configuration.

The server will listen for job submissions from clients and dispatch
them to available workers. It includes a web interface for monitoring
and managing jobs, as well as built-in Prometheus metrics.

The server uses LevelDB for persistent storage by default and supports
scheduled jobs via cron expressions.

Examples:
  # Start server with default settings
  gearhulk server

  # Start server on specific address
  gearhulk server --addr 0.0.0.0:4730

  # Start server with custom storage directory
  gearhulk server --storage-dir /var/lib/gearhulk

  # Start server with custom web interface address
  gearhulk server --web-addr :8080

  # Start server with verbose logging
  gearhulk server --addr 0.0.0.0:4730 --verbose`,
	PersistentPreRun: func(c *cobra.Command, args []string) {
		c.Flags().VisitAll(func(flag *pflag.Flag) {
			log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
		})
	},
	Run: func(cmd *cobra.Command, args []string) {
		logs.InitLogs()
		defer logs.FlushLogs()
		defer runtime.HandleCrash()
		gearmand.NewServer(cfg).Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	
	// GNU-style flags with both short and long forms
	serverCmd.Flags().StringVarP(&cfg.ListenAddr, "addr", "a", ":4730", "listening address, such as 0.0.0.0:4730")
	serverCmd.Flags().StringVarP(&cfg.Storage, "storage-dir", "s", os.TempDir()+"/gearmand", "directory where LevelDB file is stored")
	serverCmd.Flags().StringVarP(&cfg.WebAddress, "web-addr", "w", ":3000", "server HTTP API address")
	
	// Add verbose flag for logging
	serverCmd.Flags().BoolP("verbose", "v", false, "enable verbose logging")
}
