#include <iostream>
#include <string>
#include <vector>
#include <cpprest/http_client.h>
#include <cpprest/json.h>

using namespace std;
using namespace web;
using namespace web::http;
using namespace web::http::client;

const string apiKey = "YOUR_SUBSCRIPTION_KEY";
const string endpoint = "YOUR_ENDPOINT";
const string imageURL = "https://example.com/sample.jpg";

void analyzeImage() {
    http_client client(U(endpoint));
    
    // Create the request body
    json::value requestBody;
    requestBody[U("url")] = json::value::string(U(imageURL));
    requestBody[U("visualFeatures")] = json::value::array({ U("Description"), U("Categories"), U("Color") });
    requestBody[U("details")] = json::value::array({ U("Celebrities") });

    // Create the request
    client.request(methods::POST, U("/analyze"), requestBody, 
        {
            { U("Ocp-Apim-Subscription-Key"), U(apiKey) },
            { U("Content-Type"), U("application/json") }
        })
    .then([](http_response response) {
        if (response.status_code() == status_codes::OK) {
            return response.extract_json();
        }
        throw runtime_error("Image analysis failed.");
    })
    .then([](json::value jsonResponse) {
        // Process the response
        if (jsonResponse.has_field(U("description"))) {
            cout << "Image Description:" << endl;
            auto captions = jsonResponse[U("description")][U("captions")].as_array();
            for (const auto& caption : captions) {
                cout << caption[U("text")].as_string() << " (Confidence: " 
                     << caption[U("confidence")].as_double() << ")" << endl;
            }
        }

        if (jsonResponse.has_field(U("categories"))) {
            cout << "Image Categories:" << endl;
            auto categories = jsonResponse[U("categories")].as_array();
            for (const auto& category : categories) {
                cout << category[U("name")].as_string() << " (Confidence: " 
                     << category[U("score")].as_double() << ")" << endl;
            }
        }

        if (jsonResponse.has_field(U("color"))) {
            cout << "Dominant Colors:" << endl;
            auto dominantColors = jsonResponse[U("color")][U("dominantColors")].as_array();
            for (const auto& color : dominantColors) {
                cout << color.as_string() << endl;
            }
        }

        if (jsonResponse.has_field(U("celebrities"))) {
            cout << "Celebrities:" << endl;
            auto celebrities = jsonResponse[U("celebrities")].as_array();
            for (const auto& celeb : celebrities) {
                cout << celeb[U("name")].as_string() << endl;
            }
        }
    })
    .wait(); // Wait for the async operation to complete
}

int main() {
    analyzeImage();
    return 0;
}
