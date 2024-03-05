/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"lambdacrate-cli/lib"
	"lambdacrate-cli/lib/browser"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
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

		var wg sync.WaitGroup
		loginUrl, token := browser.MakeLoginUrl(config.DashboardURl)

		fmt.Printf("Press enter or go to the following link to authenticate with Lambdacrate %s (^C to quit)", loginUrl)
		inputChan := make(chan int)

		go func() {
			defer wg.Done()
			i, _ := fmt.Fscanln(os.Stdin)
			inputChan <- i
		}()

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
		//here we poll the backend to see if there was an active login session that was validated

		log.Println("hello world this is to check if the person has successfully logged in")
		for i := 0; i < 120; i++ {
			time.Sleep(time.Second)
			//todo parameterize this please
			client := http.Client{}
			url := fmt.Sprintf("%s/api/auth/cli/verify-login", config.ApiURL)
			request, err := http.NewRequest("GET", url, nil)
			request.Header.Set("content-type", "application/json")
			query := request.URL.Query()
			query.Add("token", token)
			request.URL.RawQuery = query.Encode()
			if err != nil {
				log.Println(err)
			}
			resp, err := client.Do(request)
			if err != nil {
				log.Println(err)
				continue

			}
			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println("Failed to read response body: ", err)
					continue
				}

				token := map[string]string{}
				err = json.Unmarshal(body, &token)
				if err != nil {
					log.Fatal("failed to unmarshal response: ", err)
				}
				viper.Set("api_key", token["api_key"])
				err = viper.WriteConfig()
				if err != nil {
					log.Fatal("failed to write api key to config file: ", err)
				}
				return

			} else {

			}

		}

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
