use std::error::Error;
use std::fs::File;
use std::io::{BufRead, BufReader};
use reqwest::header::{AUTHORIZATION, CONTENT_TYPE};
use serde_json::{json, Value};

// Replace these with your own values
const CLIENT_ID: &str = "your_client_id";
const CLIENT_SECRET: &str = "your_client_secret";
const TENANT_ID: &str = "your_tenant_id";
const CSV_FILE_PATH: &str = "auto_reply_messages.csv";

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
async fn set_auto_reply_message(access_token: &str, user_principal_name: &str, external_message: &str, internal_message: &str) -> Result<(), Box<dyn Error>> {
    let url = format!("https://graph.microsoft.com/v1.0/{}/mailboxSettings/automaticRepliesSetting", user_principal_name);
    let headers = [
        (AUTHORIZATION, format!("Bearer {}", access_token)),
        (CONTENT_TYPE, "application/json".to_string()),
    ];

    let body = json!({
        "status": "Scheduled",
        "externalReplyMessage": external_message,
        "internalReplyMessage": internal_message
    });

    let client = reqwest::Client::new();
    let _response = client.patch(&url)
        .headers(headers.into_iter().collect())
        .json(&body)
        .send()
        .await?;

    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let access_token = get_access_token().await?;

    let file = File::open(CSV_FILE_PATH)?;
    let reader = BufReader::new(file);

    for line in reader.lines() {
        let line = line?;
        let mut csv_reader = csv::Reader::from_reader(line.as_bytes());

        for result in csv_reader.records() {
            let record = result?;
            let user_principal_name = record.get(0).unwrap();
            let external_message = record.get(1).unwrap();
            let internal_message = record.get(2).unwrap();

            set_auto_reply_message(&access_token, user_principal_name, external_message, internal_message).await?;
            println!("Auto-reply message set successfully for {}", user_principal_name);
        }
    }

    Ok(())
}
