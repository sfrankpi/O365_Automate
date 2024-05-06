# O365_Automate

Office 365 auto-reply application.

## Introduction

Using the requests library to make the HTTP requests to the Microsoft Graph API, and the built-in csv module to read the CSV file.

The get_access_token() function retrieves an access token by making a POST request to the Microsoft identity platform endpoint. This access token is required to authenticate with the Microsoft Graph API.

The set_auto_reply_message() function sets the auto-reply message for a specific user's mailbox by making a PATCH request to the mailboxSettings/automaticRepliesSetting endpoint of the Microsoft Graph API.

The main() function reads the CSV file, iterates over each row, and calls the set_auto_reply_message() function for each user.

Remember to replace the client_id, client_secret, tenant_id, and csv_file_path variables with your own values.

The CSV file should have the following format:

```csv
user_principal_name,external_reply_message,internal_reply_message
user1@example.com,"Thank you for your email. I am currently out of the office and will respond to your message as soon as possible.","Thank you for your email. I am currently out of the office and will respond to your message as soon as possible."
user2@example.com,"I'm away from the office until [date]. I will respond to your message when I return.","I'm away from the office until [date]. I will respond to your message when I return."
```