package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/load_balancers"
	"github.com/cloudflare/cloudflare-go/v4/option"
)

var endpointName string
var endpointAction string
var accountID string

func main() {
	fmt.Println("[*] Fetching environment variables...")
	// Verify environment variables are set.
	endpointName = os.Getenv("CF_ENDPOINT_NAME")
	if endpointName == "" {
		fmt.Println("[*] WARNING - No CF_ENDPOINT_NAME provided, program will list all endpoints.")
	} else {
		endpointAction = os.Getenv("CF_ENDPOINT_ACTION")
		endpointAction = strings.ToLower(endpointAction)
		if endpointAction == "enable" || endpointAction == "disable" || endpointAction == "get" {
			fmt.Printf("[*] Endpoint: %s // Endpoint Action: %s\r\n", endpointName, endpointAction)
		} else {
			fmt.Println("[-] ERROR - CF_ENDPOINT_NAME provided, but no valid CF_ENDPOINT_ACTION was set to 'enable', 'disable', or 'get'")
			os.Exit(1)
		}
	}

	accountID = os.Getenv("CF_ACCT_ID")
	if accountID == "" {
		fmt.Println("[-] ERROR - No CF_ACCT_ID provided.")
		os.Exit(1)
	}

	apiEmail := os.Getenv("CF_API_EMAIL")
	if apiEmail == "" {
		fmt.Println("[-] ERROR - No CF_API_EMAIL provided.")
		os.Exit(1)
	}

	apiKey := os.Getenv("CF_API_KEY")
	if apiKey == "" {
		fmt.Println("[-] ERROR - No CF_API_KEY provided.")
		os.Exit(1)
	}
	fmt.Println("[+] Environment variables fetched!")

	// Open Cloudflare client
	fmt.Printf("[*] Authenticating Cloudflare client with email %s\r\n", apiEmail)
	client := cloudflare.NewClient(
		option.WithAPIKey(apiKey),
		option.WithAPIEmail(apiEmail),
	)

	fmt.Printf("[*] Authenticating Cloudflare client with Email %s for Account ID %s\r\n", apiEmail, accountID)
	lbPools, err2 := client.LoadBalancers.Pools.List(context.Background(), load_balancers.PoolListParams{AccountID: cloudflare.String(accountID)}, cloudflare.DefaultClientOptions()...)
	if err2 != nil {
		fmt.Println("[-] ERROR - Failed to list load balancer pools in Cloudflare Zone.")
		fmt.Println(err2)
	}

	if lbPools.Result == nil {
		fmt.Println("[/] No LBs in Cloudflare Zone.")
	}

	ItteratePools(lbPools.Result)

	// Could technically check sizes of pages before attempting to verify if there are additional pages.
	for {
		nextPage, err := lbPools.GetNextPage()
		if err != nil {
			fmt.Println("[-] ERROR - Failed to list the next page of load balancers in Cloudflare Zone.")
			fmt.Println(err2)
		}

		if nextPage == nil {
			//fmt.Println("[/] No additional Load Balancer pages in Cloudflare Zone.")
			break
		} else {
			ItteratePools(nextPage.Result)
		}
	}

	// Are we supposed to review the pools flagged at the end and update?
	if len(poolsToUpdate) > 0 {
		for _, pool := range poolsToUpdate {
			UpdatePool(client, pool)
		}
	}
}

var poolsToUpdate []load_balancers.Pool

func ItteratePools(lbs []load_balancers.Pool) {
	for _, pool := range lbs {
		fmt.Printf("\tLoad Balancer Pool ID: %s || LB Pool Name: %s // Enabled: %v\r\n", pool.ID, pool.Name, pool.Enabled)
		poolChanged := false
		for k, origin := range pool.Origins {
			if endpointName == "" {
				fmt.Printf("\tLoad Balancer Pool Origin: %s || Address: %s / VNet: %s // Enabled: %v\r\n", origin.Name, origin.Address, origin.VirtualNetworkID, origin.Enabled)
			} else if endpointName == origin.Name {
				fmt.Printf("\t** Load Balancer Pool Origin: %s || Address: %s / VNet: %s // Enabled: %v\r\n", origin.Name, origin.Address, origin.VirtualNetworkID, origin.Enabled)
			}

			// Trigger endpoint action?
			if origin.Name == endpointName {
				if endpointAction == "enable" {
					pool.Origins[k].Enabled = true
					poolChanged = true
				} else if endpointAction == "disable" {
					pool.Origins[k].Enabled = false
					poolChanged = true
				} else {

				}

				// Notify in terminal of update?
				if endpointAction == "enable" || endpointAction == "disable" {
					fmt.Printf("\t*UPDATED*\t** Load Balancer Pool Origin: %s || Address: %s / VNet: %s // Enabled: %v\r\n", origin.Name, origin.Address, origin.VirtualNetworkID, origin.Enabled)
				}
			}
		}

		// Add pool to be updaet if updated.
		if poolChanged {
			poolsToUpdate = append(poolsToUpdate, pool)
		}
	}
}

func UpdatePool(client *cloudflare.Client, pool load_balancers.Pool) {
	fmt.Printf("[*] Updating Load Balancer Pool ID: %s || Pool Name: %s\r\n", pool.ID, pool.Name)

	originParam := []load_balancers.OriginParam{}
	for _, origin := range pool.Origins {
		originParam = append(originParam, load_balancers.OriginParam{
			Name:    cloudflare.F(origin.Name),
			Address: cloudflare.F(origin.Address),
			Enabled: cloudflare.F(origin.Enabled),
		})
	}

	fmt.Printf("[*] Sending API Request... \r\n")
	res, err1 := client.LoadBalancers.Pools.Update(context.Background(), pool.ID,
		load_balancers.PoolUpdateParams{AccountID: cloudflare.String(accountID),
			Name: cloudflare.String(pool.Name), Origins: cloudflare.F(originParam)},
		cloudflare.DefaultClientOptions()...)

	if err1 != nil {
		fmt.Printf("[-] ERROR: Failed to update Load Balancer Pool ID: %s || Pool Name: %s\r\n", pool.ID, pool.Name)
		fmt.Println(err1)
	} else {
		fmt.Printf("[*] API Request Complete.\r\n\tUpdated Load Balancer Pool ID: %s || Pool Name: %s\r\n\t** PLEASE NOTE: Changes may take a brief period of time to be visible in the dashboard.\r\n", res.ID, res.Name)
	}
}
