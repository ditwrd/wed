/*
Copyright © 2025 Aditya Wardianto <hi@ditwrd.dev>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"github.com/ditwrd/wed/internal/db"
	"github.com/ditwrd/wed/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		db.Module,
		server.Module,
		server.PageModule,
	).Run()
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	serveCmd.PersistentFlags().String("bind", "", "bind address (overrides app.bind)")
	serveCmd.PersistentFlags().Int("port", 0, "port (overrides app.port)")
	serveCmd.PersistentFlags().Bool("dev", false, "development mode (overrides app.dev)")

	_ = viper.BindPFlag("app.bind", serveCmd.PersistentFlags().Lookup("bind"))
	_ = viper.BindPFlag("app.port", serveCmd.PersistentFlags().Lookup("port"))
	_ = viper.BindPFlag("app.dev", serveCmd.PersistentFlags().Lookup("dev"))
}
