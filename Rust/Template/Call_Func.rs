use serde_json::{Value};
use reqwest;
use std::collections::HashMap;

async fn call_azure_cognitive_service(image_url: &str) -> Result<(), Box<dyn std::error::Error>> {
    let subscription_key = "YOUR_SUBSCRIPTION_KEY";
    let endpoint = "YOUR_ENDPOINT";
    let url = format!("{}/vision/v3.1/analyze?visualFeatures=Categories,Description,Color&details=Celebrities&language=en", endpoint);
    let client = reqwest::Client::new();

    let mut headers = reqwest::header::HeaderMap::new();
    headers.insert("Ocp-Apim-Subscription-Key", subscription_key.parse().unwrap());
    headers.insert(reqwest::header::CONTENT_TYPE, "application/json".parse().unwrap());

    let body = json!({ "url": image_url });

    let response = client.post(&url)
        .headers(headers)
        .json(&body)
        .send()
        .await?;

    let response_text = response.text().await?;
    let json_response: Value = serde_json::from_str(&response_text)?;

    // Parse and display the response
    let description = &json_response["description"]["captions"][0]["text"];
    let categories = &json_response["categories"];
    let colors = &json_response["color"]["dominantColors"];
    let celebrities = &json_response["celebrities"];

    println!("Image Description: {}", description);
    println!("Image Categories: {}", categories);
    println!("Dominant Colors: {}", colors);
    println!("Celebrities: {}", celebrities);

    Ok(())
}

#[tokio::main]
async fn main() {
    let image_url = "https://example.com/sample.jpg";

    match call_azure_cognitive_service(image_url).await {
        Ok(_) => println!("Image analyzed successfully."),
        Err(err) => eprintln!("Error: {}", err),
    }
}