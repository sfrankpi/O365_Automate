import requests
import json

# Replace these with your own values
client_id = 'your_client_id'
client_secret = 'your_client_secret'
tenant_id = 'your_tenant_id'

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
def set_auto_reply_message(access_token, message):
    url = 'https://graph.microsoft.com/v1.0/me/mailboxSettings/automaticRepliesSetting'
    headers = {
        'Authorization': f'Bearer {access_token}',
        'Content-Type': 'application/json'
    }
    data = {
        'status': 'Scheduled',
        'externalReplyMessage': message,
        'internalReplyMessage': message
    }
    response = requests.patch(url, headers=headers, data=json.dumps(data))
    response.raise_for_status()
    return response.json()

# Example usage
def main():
    try:
        access_token = get_access_token()
        auto_reply_message = 'Thank you for your email. I am currently out of the office and will respond to your message as soon as possible.'
        set_auto_reply_message(access_token, auto_reply_message)
        print('Auto-reply message set successfully!')
    except requests.exceptions.RequestException as e:
        print(f'Error setting auto-reply message: {e}')

if __name__ == '__main__':
    main()
