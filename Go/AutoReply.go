package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Replace these with your own values
const clientID = "your_client_id"
const clientSecret = "your_client_secret"
const tenantID = "your_tenant_id"

// Function to get an access token
func getAccessToken() (string, error) {
	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	data := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"scope":         "https://graph.microsoft.com/.default",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", err
	}

	return respData["access_token"].(string), nil
}

// Function to set the auto-reply message
func setAutoReplyMessage(accessToken, message string) error {
	url := "https://graph.microsoft.com/v1.0/me/mailboxSettings/automaticRepliesSetting"
	data := map[string]interface{}{
		"status":               "Scheduled",
		"externalReplyMessage": message,
		"internalReplyMessage": message,
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

	autoReplyMessage := "Thank you for your email. I am currently out of the office and will respond to your message as soon as possible."
	err = setAutoReplyMessage(accessToken, autoReplyMessage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting auto-reply message: %v\n", err)
		return
	}

	fmt.Println("Auto-reply message set successfully!")
}
