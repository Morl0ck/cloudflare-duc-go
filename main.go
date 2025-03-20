package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

func main() {
	// Load environment variables
	apiToken := os.Getenv("CF_API_TOKEN")
	zoneName := os.Getenv("CF_ZONE_NAME")
	recordName := os.Getenv("CF_RECORD_NAME")
	interval := os.Getenv("UPDATE_INTERVAL") // in minutes

	if apiToken == "" || zoneName == "" || recordName == "" {
		log.Fatal("CF_API_TOKEN, CF_ZONE_NAME, and CF_RECORD_NAME must be set")
	}

	// Set default interval if not set
	updateInterval := 5 * time.Minute
	if interval != "" {
		if i, err := time.ParseDuration(interval + "m"); err == nil {
			updateInterval = i
		} else {
			log.Printf("Invalid UPDATE_INTERVAL '%s', using default 5 minutes", interval)
		}
	}

	// Initialize Cloudflare API
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatalf("Failed to create Cloudflare API client: %v", err)
	}

	// Get Zone ID
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		log.Fatalf("Failed to get Zone ID for %s: %v", zoneName, err)
	}

	// Get DNS Record ID
	recordID, recordIP, err := getDNSRecord(api, zoneID, recordName)
	if err != nil {
		log.Fatalf("Failed to get DNS record: %v", err)
	}

	log.Printf("Current DNS record %s points to %s", recordName, recordIP)

	// Start ticker
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentIP, err := getPublicIP()
			if err != nil {
				log.Printf("Error fetching public IP: %v", err)
				continue
			}

			if !isValidIP(currentIP) {
				log.Printf("Invalid IP address: '%s'", currentIP)
				continue
			}
	
			log.Printf("Current public IP: %s", currentIP)

			if currentIP != recordIP {
				log.Printf("IP has changed from %s to %s. Updating DNS record...", recordIP, currentIP)
				err = updateDNSRecord(api, zoneID, recordID, currentIP)
				if err != nil {
					log.Printf("Error updating DNS record: %v", err)
				} else {
					log.Printf("Successfully updated DNS record to %s", currentIP)
					recordIP = currentIP
				}
			} else {
				log.Println("IP address has not changed. No update needed.")
			}
		}
	}
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var ip string
	_, err = fmt.Fscanf(resp.Body, "%s", &ip)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func getDNSRecord(api *cloudflare.API, zoneID, recordName string) (string, string, error) {
	records, _, err := api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
		Name: recordName,
		Type: "A",
	})

	if err != nil {
		return "", "", err
	}

	if len(records) == 0 {
		return "", "", fmt.Errorf("no A record found for %s", recordName)
	}

	return records[0].ID, records[0].Content, nil
}

func updateDNSRecord(api *cloudflare.API, zoneID, recordID, newIP string) error {
	_, err := api.UpdateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
		ID:      recordID,
		Content: newIP,
	})
	return err
}

func isValidIP(ip string) bool {
    return net.ParseIP(ip) != nil
}
