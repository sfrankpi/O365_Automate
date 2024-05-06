package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Replace these with your own values
const clientID = "your_client_id"
const clientSecret = "your_client_secret"
const tenantID = "your_tenant_id"
const csvFilePath = "auto_reply_messages.csv"

type autoReplyMessage struct {
	UserPrincipalName string
	ExternalMessage   string
	InternalMessage   string
}

// Function to get an access token
func getAccessToken() (string, error) {
	// ... (same as previous example)
}

// Function to set the auto-reply message
func setAutoReplyMessage(accessToken string, message *autoReplyMessage) error {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/%s/mailboxSettings/automaticRepliesSetting", message.UserPrincipalName)
	data := map[string]interface{}{
		"status":               "Scheduled",
		"externalReplyMessage": message.ExternalMessage,
		"internalReplyMessage": message.InternalMessage,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error setting auto-reply message: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	accessToken, err := getAccessToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting access token: %v\n", err)
		return
	}

	file, err := os.Open(csvFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening CSV file: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading CSV file: %v\n", err)
			return
		}

		message := &autoReplyMessage{
			UserPrincipalName: row[0],
			ExternalMessage:   row[1],
			InternalMessage:   row[2],
		}

		err = setAutoReplyMessage(accessToken, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error setting auto-reply message for %s: %v\n", message.UserPrincipalName, err)
			continue
		}

		fmt.Printf("Auto-reply message set successfully for %s\n", message.UserPrincipalName)
	}
}
