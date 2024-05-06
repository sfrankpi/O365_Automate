import requests
import csv
import json

# Replace these with your own values
client_id = 'your_client_id'
client_secret = 'your_client_secret'
tenant_id = 'your_tenant_id'
csv_file_path = 'auto_reply_messages.csv'

# Function to get an access token
def get_access_token():
    url = f'https://login.microsoftonline.com/{tenant_id}/oauth2/v2.0/token'
    data = {
        'grant_type': 'client_credentials',
        'client_id': client_id,
        'client_secret': client_secret,
        'scope': 'https://graph.microsoft.com/.default'
    }
    response = requests.post(url, data=data)
    response.raise_for_status()
    return response.json()['access_token']

# Function to set the auto-reply message
def set_auto_reply_message(access_token, user_principal_name, external_message, internal_message):
    url = f'https://graph.microsoft.com/v1.0/{user_principal_name}/mailboxSettings/automaticRepliesSetting'
    headers = {
        'Authorization': f'Bearer {access_token}',
        'Content-Type': 'application/json'
    }
    data = {
        'status': 'Scheduled',
        'externalReplyMessage': external_message,
        'internalReplyMessage': internal_message
    }
    response = requests.patch(url, headers=headers, data=json.dumps(data))
    response.raise_for_status()
    return response.json()

def main():
    try:
        access_token = get_access_token()

        with open(csv_file_path, 'r') as csv_file:
            reader = csv.DictReader(csv_file)
            for row in reader:
                user_principal_name = row['user_principal_name']
                external_message = row['external_reply_message']
                internal_message = row['internal_reply_message']
                set_auto_reply_message(access_token, user_principal_name, external_message, internal_message)
                print(f'Auto-reply message set successfully for {user_principal_name}')

    except requests.exceptions.RequestException as e:
        print(f'Error setting auto-reply messages: {e}')

if __name__ == '__main__':
    main()
