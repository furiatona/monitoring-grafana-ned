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
	_ = godotenv.Load() // optional: load .env if present
}

func main() {
	// Override with env vars only if they're non-empty
	vars := map[*string]string{
		&consul.ConsulAddress: "EMBED_CONSUL_ADDRESS",
		&consul.ConsulToken:   "EMBED_CONSUL_TOKEN",
		&consul.ConsulKVKey:   "EMBED_CONSUL_KV_KEY",
	}
	for ptr, env := range vars {
		if val := os.Getenv(env); val != "" {
			*ptr = val
		}
	}

	// DEBUG: Show loaded config (optional)
	/*
		fmt.Println("ConsulAddress:", consul.ConsulAddress)
		fmt.Println("ConsulToken:", consul.ConsulToken)
		fmt.Println("ConsulKVKey:", consul.ConsulKVKey)
	*/

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
