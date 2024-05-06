// the necessary Office 365 API credentials
const clientId = 'your_client_id';
const clientSecret = 'your_client_secret';
const tenantId = 'your_tenant_id';

// Import necessary libraries
const axios = require('axios');

// Function to authenticate with the Office 365 API
async function getAccessToken() {
  const response = await axios.post(`https://login.microsoftonline.com/${tenantId}/oauth2/v2.0/token`, {
    grant_type: 'client_credentials',
    client_id: clientId,
    client_secret: clientSecret,
    scope: 'https://graph.microsoft.com/.default'
  });

  return response.data.access_token;
}

// Function to set the auto-reply message
async function setAutoReplyMessage(accessToken, message) {
  const response = await axios.patch(
    'https://graph.microsoft.com/v1.0/me/mailboxSettings/automaticRepliesSetting',
    {
      status: 'Scheduled',
      externalReplyMessage: message,
      internalReplyMessage: message
    },
    {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      }
    }
  );

  return response.data;
}

// Example usage
async function main() {
  try {
    const accessToken = await getAccessToken();
    const autoReplyMessage = 'Thank you for your email. I am currently out of the office and will respond to your message as soon as possible.';
    await setAutoReplyMessage(accessToken, autoReplyMessage);
    console.log('Auto-reply message set successfully!');
  } catch (error) {
    console.error('Error setting auto-reply message:', error);
  }
}

main();
