#include <iostream>
#include <fstream>
#include <string>
#include <vector>
#include <cpprest/http_client.h>
#include <cpprest/json.h>
#include <cpprest/uri.h>
#include <sstream>

using namespace std;
using namespace web;
using namespace web::http;
using namespace web::http::client;

const string client_id = "your_client_id";
const string client_secret = "your_client_secret";
const string tenant_id = "your_tenant_id";
const string csv_file_path = "auto_reply_messages.csv";

string get_access_token() {
    http_client client(U("https://login.microsoftonline.com/" + tenant_id + "/oauth2/v2.0/token"));
    
    json::value data;
    data[U("grant_type")] = json::value::string(U("client_credentials"));
    data[U("client_id")] = json::value::string(U(client_id));
    data[U("client_secret")] = json::value::string(U(client_secret));
    data[U("scope")] = json::value::string(U("https://graph.microsoft.com/.default"));

    client.request(methods::POST, U(""), data, { { U("Content-Type"), U("application/x-www-form-urlencoded") } })
        .then([](http_response response) {
            if (response.status_code() == status_codes::OK) {
                return response.extract_json();
            }
            throw runtime_error("Failed to get access token.");
        })
        .then([](json::value jsonResponse) {
            return jsonResponse[U("access_token")].as_string();
        })
        .wait();
}

void set_auto_reply_message(const string& access_token, const string& user_principal_name, const string& external_message, const string& internal_message) {
    http_client client(U("https://graph.microsoft.com/v1.0/" + user_principal_name + "/mailboxSettings/automaticRepliesSetting"));
    
    json::value data;
    data[U("status")] = json::value::string(U("Scheduled"));
    data[U("externalReplyMessage")] = json::value::string(U(external_message));
    data[U("internalReplyMessage")] = json::value::string(U(internal_message));

    client.request(methods::PATCH, U(""), data, { { U("Authorization"), U("Bearer " + access_token) }, { U("Content-Type"), U("application/json") } })
        .then([](http_response response) {
            if (response.status_code() != status_codes::OK) {
                throw runtime_error("Failed to set auto-reply message.");
            }
        })
        .wait();
}

void read_csv_and_set_replies() {
    try {
        string access_token = get_access_token();

        ifstream csv_file(csv_file_path);
        string line;
        getline(csv_file, line); // Read header line

        while (getline(csv_file, line)) {
            stringstream ss(line);
            string user_principal_name, external_message, internal_message;

            getline(ss, user_principal_name, ',');
            getline(ss, external_message, ',');
            getline(ss, internal_message, ',');

            set_auto_reply_message(access_token, user_principal_name, external_message, internal_message);
            cout << "Auto-reply message set successfully for " << user_principal_name << endl;
        }
    } catch (const exception& e) {
        cout << "Error setting auto-reply messages: " << e.what() << endl;
    }
}

int main() {
    read_csv_and_set_replies();
    return 0;
}
