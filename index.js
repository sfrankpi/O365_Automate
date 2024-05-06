const axios = require('axios');

// Replace these with your own values
const clientId = 'your_client_id';
const clientSecret = 'your_client_secret';
const tenantId = 'your_tenant_id';
const csvFilePath = 'auto_reply_messages.csv';

// Function to get an access token
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
async function setAutoReplyMessage(accessToken, userPrincipalName, externalMessage, internalMessage) {
  const response = await axios.patch(
    `https://graph.microsoft.com/v1.0/${userPrincipalName}/mailboxSettings/automaticRepliesSetting`,
    {
      status: 'Scheduled',
      externalReplyMessage: externalMessage,
      internalReplyMessage: internalMessage
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

async function main() {
  try {
    const accessToken = await getAccessToken();

    const fs = require('fs');
    const csv = require('csv-parser');

    fs.createReadStream(csvFilePath)
      .pipe(csv())
      .on('data', async (row) => {
        const { user_principal_name, external_reply_message, internal_reply_message } = row;
        await setAutoReplyMessage(accessToken, user_principal_name, external_reply_message, internal_reply_message);
        console.log(`Auto-reply message set successfully for ${user_principal_name}`);
      })
      .on('end', () => {
        console.log('All auto-reply messages set successfully!');
      });
  } catch (error) {
    console.error('Error setting auto-reply messages:', error);
  }
}

main();
