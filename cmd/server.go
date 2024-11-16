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
	Use: "server",
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
	serverCmd.Flags().StringVar(&cfg.ListenAddr, "addr", ":4730", "listening on, such as 0.0.0.0:4730")
	serverCmd.Flags().StringVar(&cfg.Storage, "storage-dir", os.TempDir()+"/gearmand", "Directory where LevelDB file is stored.")
	serverCmd.Flags().StringVar(&cfg.WebAddress, "web.addr", ":3000", "Server HTTP api Address")
}
