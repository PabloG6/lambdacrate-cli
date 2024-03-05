/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"lambdacrate-cli/lib"
	"lambdacrate-cli/lib/auth"
	"lambdacrate-cli/lib/browser"
	"log"
	"os"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := lib.LoadConfig()
		_ = context.Background()

		if err != nil {
			log.Fatal("unable to load config file with error: ", err)
		}

		loginUrl, token := browser.MakeLoginUrl(config.DashboardURl)

		fmt.Printf("Press enter or go to the following link to authenticate with Lambdacrate %s (^C to quit)", loginUrl)
		inputChan := make(chan int)
		pollResponse := make(chan auth.PollConfirmAuthResponse)
		//start polling the api for the confirmation response.
		go func() {

			select {
			case <-inputChan:
				{
					err = browser.Open(loginUrl)
					if err != nil {
						fmt.Printf("Unable to open %s. Please copy or click on the link", loginUrl)
						fmt.Printf("\n\n\n")

					}

					return
				}
			}

		}()
		go auth.AsyncPollConfirmAuth(config, auth.PollConfirmAuthRequest{
			Token: token,
		}, pollResponse)
		go func() {
			i, _ := fmt.Fscanln(os.Stdin)
			inputChan <- i
		}()

		//here we poll the backend to see if there was an active login session that was validated

	},
}

func openBrowser() {

}
func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
