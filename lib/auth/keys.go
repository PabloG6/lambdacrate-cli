package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"lambdacrate-cli/lib"
	"log"
	"net/http"
	"time"
)

type PollConfirmAuthRequest struct {
	Token string
}

type PollConfirmAuthResponse struct {
	ApiKey string `json:"api_key"`
}

func AsyncPollConfirmAuth(config lib.Config, pollRequest PollConfirmAuthRequest, outChan chan PollConfirmAuthResponse) {
	defer close(outChan)
	for i := 0; i < 120; i++ {
		time.Sleep(time.Second)
		//todo parameterize this please
		client := http.Client{}
		url := fmt.Sprintf("%s/api/auth/cli/verify-login", config.ApiURL)
		request, err := http.NewRequest("GET", url, nil)
		request.Header.Set("content-type", "application/json")
		query := request.URL.Query()
		query.Add("token", pollRequest.Token)
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
			outChan <- PollConfirmAuthResponse{ApiKey: token["api_key"]}

			return

		}

	}

}
