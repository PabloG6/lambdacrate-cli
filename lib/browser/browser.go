package browser

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

// borrowed from
var execCommand = exec.Command

func MakeLoginUrl(baseUrl string) (string, string) {
	numBytes := 32

	// Generate random bytes
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("unable to generate random login token. rand package unavailable?: ", err)
		return "", ""
	}

	// Encode random bytes to base64
	token := base64.RawURLEncoding.EncodeToString(randomBytes)

	return fmt.Sprintf("%s/login/cli/confirm_auth?t=%s", baseUrl, token), token
}
func Open(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = execCommand("xdg-open", url).Start()
	case "windows":
		err = execCommand("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = execCommand("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return err
	}

	return nil
}
