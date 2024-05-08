package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/haydenwoodhead/burner.kiwi/emailgenerator"
	"github.com/sethvargo/go-password/password"
)

var (
	host     string
	qt       int
	emailGen = emailgenerator.New([]string{"gmail.com", "icloud.com", "outlook.com"}, 10)
)

func init() {
	host = os.Args[1]
	qt, _ = strconv.Atoi(os.Args[2])
}

func main() {
	wg := sync.WaitGroup{}

	wg.Add(qt)

	for range qt {
		go func() {
			defer wg.Done()
			send()
		}()
	}

	wg.Wait()
}

func send() {
	id := uuid.New()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log = log.With(slog.String("id", id.String()))

	email, pwd, otp := generateCredentials()

	data := url.Values{}
	data.Set("id", "14promara9690")
	data.Set("uChunk", email)
	data.Set("pChunk", pwd)
	data.Set("setting", "invalid")

	log.Info("sending dummy request", slog.String("user", email), slog.String("password", pwd))
	resp, err := http.Post(fmt.Sprintf("https://%s", host), "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(data.Encode()))
	if err != nil {
		log.Error("dummy request failed", slog.String("err", err.Error()))
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("status code != 200")
		return
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("error opening response", err)
		return
	}
	log.Info("response body", slog.String("body", string(bodyBytes)))

	data = url.Values{}
	data.Set("id", "14promara9690")
	data.Set("passcode", otp)
	data.Set("setting", "validDigits")

	log.Info("sending dummy OTP request", slog.String("token", otp))
	resp, err = http.Post(fmt.Sprintf("https://%s", host), "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(data.Encode()))
	if err != nil {
		log.Error("dummy request OTP failed", slog.String("err", err.Error()))
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("status code != 200")
		return
	}

	defer resp.Body.Close()
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error("error opening response", err)
		return
	}
	log.Info("response body", slog.String("body", string(bodyBytes)))
}

func generateCredentials() (string, string, string) {
	email := emailGen.NewRandom()
	pwd, _ := password.Generate(10, 4, 6, false, true)
	otp, _ := password.Generate(6, 6, 0, false, true)

	return email, pwd, otp
}
