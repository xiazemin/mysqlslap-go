/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"context"
	"fmt"
	"mysqlslap/db"
	"mysqlslap/run"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var host string
var port int64
var user string
var password string
var database string
var sql string
var concurrency int
var iteration int
var timeout int64

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mysqlslap",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mysqlslap -H", host,
			" -p", port,
			" -u", user,
			" -P", password,
			" -d", database,
			" -q", sql,
			" -c", concurrency,
			" -i", iteration, "\n", "start...")

		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()

		db, err := db.GetClient(ctx, host, port, user, password, database, concurrency, timeout)
		if err != nil {
			panic(err)
		}
		run.Run(ctx, db, sql, concurrency, iteration)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "mysql host name,default 127.0.0.1")

	rootCmd.PersistentFlags().Int64VarP(&port, "port", "P", 3306, "mysql port,default 3306")

	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "root", "user name ,default root")

	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "12345678", "password,default 12345678")

	rootCmd.PersistentFlags().StringVarP(&database, "database", "d", "mysql", "database name,default mysql")

	rootCmd.PersistentFlags().StringVarP(&sql, "query", "q", "show tables;", "the query you want to run,default show tables")

	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 30, "concurrency,default 30")

	rootCmd.PersistentFlags().IntVarP(&iteration, "iteration", "i", 200, "number of iterations ,default  200")

	rootCmd.PersistentFlags().Int64VarP(&timeout, "timeout", "t", 3600, "db connect time out (second),default 3600")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mysqlslap.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mysqlslap" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mysqlslap")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
