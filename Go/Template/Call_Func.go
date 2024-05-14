// go get github.com/Azure/azure-sdk-for-go
package main

import (
    "context"
    "fmt"
    "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/computervision"
    "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/computervision/computervision"
    "github.com/Azure/go-autorest/autorest"
)

func main() {
    // Set up the credentials
    apiKey := "YOUR_SUBSCRIPTION_KEY"
    endpoint := "YOUR_ENDPOINT"
    creds := NewCognitiveServicesCredentials(apiKey)
    authorizer := autorest.NewCognitiveServicesAuthorizer(creds)

    // Create a new client
    client := computervision.New(endpoint)
    client.Authorizer = authorizer

    // Analyze an image
    imageURL := "https://example.com/sample.jpg"
    features := []computervision.VisualFeatureTypes{computervision.Description, computervision.Categories, computervision.Color}
    info, err := client.AnalyzeImage(context.Background(), imageURL, features, []computervision.Details{computervision.Celebrities})

    if err != nil {
        fmt.Println("Image analysis failed.")
        return
    }

    // Display the analysis results
    if info.Description != nil {
        fmt.Println("Image Description:")
        for _, caption := range *info.Description.Captions {
            fmt.Printf("%s (Confidence: %.2f)\n", *caption.Text, *caption.Confidence)
        }
    }

    if info.Categories != nil {
        fmt.Println("Image Categories:")
        for _, category := range *info.Categories {
            fmt.Printf("%s (Confidence: %.2f)\n", *category.Name, *category.Score)
        }
    }

    if info.Color != nil {
        fmt.Println("Dominant Colors:")
        for _, color := range *info.Color.DominantColors {
            fmt.Println(*color)
        }
    }

    if info.Celebrities != nil {
        fmt.Println("Celebrities:")
        for _, celeb := range *info.Celebrities {
            fmt.Println(*celeb.Name)
        }
    }
}