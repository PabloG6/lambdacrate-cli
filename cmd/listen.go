/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"lambdacrate-cli/lib"
	"lambdacrate-cli/lib/proxy"
	"log"
	"sync"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup

		fmt.Println("listen called")
		port, err := cmd.Flags().GetString("forward-to")

		if err != nil {
			log.Fatal("No port passed to command. Terminating...")
		}

		app, err := cmd.Flags().GetString("app")

		if err != nil {
			log.Fatal("No app flag passed to command. Terminating...")
		}

		client, err := lib.NewClient(port, &wg)

		proxy := proxy.NewProxy(app, client)

		go func() {
			err = proxy.Run()
			if err != nil {
				msg := fmt.Sprintf("Unable to establish a connection with the remote server for app %s", app)
				log.Fatal(msg)
			} else {
				log.Println("Successfully established connection to remote server, proxy is now active. ")
			}
		}()

		//this is the main loop of the application
		for {
			select {
			case <-proxy.CloseChan:
				{
					log.Fatal("Terminating connection...")
				}
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	// Here you will define your flags and configuration settings.
	listenCmd.PersistentFlags().String("app", "", "The app you'd like to make a proxy to")
	listenCmd.PersistentFlags().String("forward-to", "", "Local port your app is listening on.")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//helper functions

func getForwardUrl(port string) string {
	return fmt.Sprintf("http://localhost:%s", port)
}
