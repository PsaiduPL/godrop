package client

import (
	"bytes"
	"errors"
	"fmt"
	"godrop/internal/common"
	"godrop/internal/common/colored"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// prevalidated
type ClientOptions struct {
	From     *string // required
	Password *string // default nil
	Port     uint
}

// postvalidated
type ClientConfig struct {
	from           string
	hashedPassword []byte
	port           uint
}

func GetFile(clientOptions ClientOptions) error {
	clientConfig, err := clientOptions.validate()
	if err != nil {
		return err
	}

	client := http.Client{

		Timeout: 30 * time.Second,
	}
	url, err := url.Parse(fmt.Sprintf("http://%s:%d/api/v1/share", clientConfig.from, clientConfig.port))
	if err != nil {
		slog.Warn("Invalid url", "error", err)
		return err
	}
	header := prepareHeader(clientConfig)
	req, err := http.NewRequest(http.MethodGet, url.String(), bytes.NewReader(nil))

	for k, v := range header {
		req.Header.Add(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		slog.Warn("Error while making request", "error", err)
		return err
	}
	return handleResponse(res)

}

func (clientOptions *ClientOptions) validate() (*ClientConfig, error) {
	var password = common.DefaultPassword
	var from string
	if clientOptions.From == nil || strings.TrimSpace(*clientOptions.From) == "" {
		return nil, fmt.Errorf(colored.BuildColoredString("Invalid <RED>IP<RED> cannot be empty"))
	}
	from = *clientOptions.From

	if clientOptions.Password != nil {
		password = *clientOptions.Password
	}
	passHashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &ClientConfig{
		hashedPassword: passHashed,
		from:           from,
		port:           clientOptions.Port,
	}, nil
}

func prepareHeader(clientConfig *ClientConfig) map[string]string {
	return map[string]string{
		common.HeaderName: string(clientConfig.hashedPassword),
	}
}

func handleResponse(response *http.Response) error {

	defer response.Body.Close()
	statusCode := response.StatusCode
	if statusCode >= 400 && statusCode < 500 {
		return fmt.Errorf("Error while making request to server check password")
	}

	header := response.Header
	fileName := header.Get("X-FileName")
	if fileExists(fileName) {
		return fmt.Errorf("File with specified name exists rename file")
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	body := response.Body
	_, err = io.Copy(file, body)
	if err == nil {
		colored.PrintColoredWithTags("<GREEN>File was succesfuly downloaded<GREEN> file <RED>%s<RED>\n", file.Name())
	}
	return err
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrExist) {
		return true
	}
	return false
}
