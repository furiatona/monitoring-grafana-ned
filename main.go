package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"

	ui "ned/internal/cli"
	consul "ned/internal/conn"
	dcpatterns "ned/internal/discovery"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, proceeding without it.")
	}
}

func main() {
	// Set Consul config from env if not already set
	vars := map[*string]string{
		&consul.ConsulAddress: "CONSUL_ADDRESS",
		&consul.ConsulToken:   "CONSUL_TOKEN",
		&consul.ConsulKVKey:   "CONSUL_KV_KEY",
	}
	for ptr, env := range vars {
		if *ptr == "" {
			*ptr = os.Getenv(env)
		}
	}

	if consul.ConsulAddress == "" || consul.ConsulKVKey == "" {
		fmt.Println("Error: CONSUL_ADDRESS and CONSUL_KV_KEY must be set.")
		os.Exit(1)
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		os.Exit(1)
	}
	instanceLabel := strings.Split(hostname, ".")[0]

	patterns, err := consul.FetchDCPatterns()
	if err != nil {
		fmt.Println("Error fetching DC patterns:", err)
		os.Exit(1)
	}

	recommended := dcpatterns.GetRecommendedDC(instanceLabel, patterns)
	uniqueDCs := dcpatterns.GetUniqueDCLocations(patterns)
	dcLocation := ui.SelectDC(recommended, uniqueDCs)

	fmt.Printf("\033[0;32mDC_LOCATION=%s\033[0m\n", dcLocation)
}
