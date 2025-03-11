package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Noviiich/vpn-config-generator/config"
)

func main() {
	cfg := config.Load()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.timeweb.cloud/api/v1/servers/4383899", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.ServerToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

}
