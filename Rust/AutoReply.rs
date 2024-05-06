use std::error::Error;
use reqwest::header::{AUTHORIZATION, CONTENT_TYPE};
use serde_json::{json, Value};

// Replace these with your own values
const CLIENT_ID: &str = "your_client_id";
const CLIENT_SECRET: &str = "your_client_secret";
const TENANT_ID: &str = "your_tenant_id";

// Function to get an access token
async fn get_access_token() -> Result<String, Box<dyn Error>> {
    let url = format!("https://login.microsoftonline.com/{}/oauth2/v2.0/token", TENANT_ID);
    let params = [
        ("grant_type", "client_credentials"),
        ("client_id", CLIENT_ID),
        ("client_secret", CLIENT_SECRET),
        ("scope", "https://graph.microsoft.com/.default"),
    ];

    let client = reqwest::Client::new();
    let response = client.post(&url)
        .form(&params)
        .send()
        .await?;

    let json_response: Value = response.json().await?;
    Ok(json_response["access_token"].as_str().unwrap().to_string())
}

// Function to set the auto-reply message
async fn set_auto_reply_message(access_token: &str, message: &str) -> Result<(), Box<dyn Error>> {
    let url = "https://graph.microsoft.com/v1.0/me/mailboxSettings/automaticRepliesSetting";
    let headers = [
        (AUTHORIZATION, format!("Bearer {}", access_token)),
        (CONTENT_TYPE, "application/json".to_string()),
    ];

    let body = json!({
        "status": "Scheduled",
        "externalReplyMessage": message,
        "internalReplyMessage": message
    });

    let client = reqwest::Client::new();
    let _response = client.patch(url)
        .headers(headers.into_iter().collect())
        .json(&body)
        .send()
        .await?;

    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let access_token = get_access_token().await?;
    let auto_reply_message = "Thank you for your email. I am currently out of the office and will respond to your message as soon as possible.";
    set_auto_reply_message(&access_token, auto_reply_message).await?;
    println!("Auto-reply message set successfully!");
    Ok(())
}
