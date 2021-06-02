## Prerequisites

- [Go](https://golang.org/doc/install)
- [ngrok](https://ngrok.com/download)

    ```bash
    # start a HTTP tunnel on port 8080
    ./ngrok http 8080
    ```

**Attention!** Ngrok tunneling should not be used for production applications and was used here only for demonstration ease of use.

## Creating a LiveChat app

1. Sign in to [Developer Console](https://developers.livechatinc.com/console/).
2. Go to the [Apps](https://developers.livechatinc.com/console/apps/) section.
3. Create a new app. You can use the Blank template and add the Authorization building block manually or use the Server-side webhook app template.
4. Configure the Authorization building block.
	1. Choose server-side as the Client type.
	2. Add `$NGROK_PUBLIC_URL/oauth` to the Redirect URI whitelist.
	3. Copy Client Id, Client Secret, and Redirect URI. You will need them when creating a configuration file for the app.
	4. Add `webhooks--all:rw` and `chats--all:rw` to the list of requested scopes.
	5. Use `$NGROK_PUBLIC_URL/oauth` as the Direct installation URL in the Marketplace authorization flow settings section.

**Attention!** For production applications, you can skip point 4.5.
Having troubles? Visit docs dedicated page for [creating LiveChat apps](https://developers.livechatinc.com/docs/getting-started/guides/#creating-livechat-apps).

## Running this example

- Clone the repository

    ```bash
    git clone https://github.com/livechat/lc-sdk-go/v4.git
    ```

- In a terminal, navigate to `lc-sdk-go/examples/echo`

    ```bash
    cd lc-sdk-go/examples/echo
    ```

- Create a config file from the example

    ```bash
    cp config.example.json config.json # for Unix and Unix-like operating systems
    copy config.example.json config.json # for Windows
    ```

    Use values saved when [creating a LiveChat app](#creating-a-livechat-apps) to fill `config.json`.
    Assuming that `$NGROK_PUBLIC_URL` is equal to `https://3c6129e3.ngrok.io`, sample config file looks like:

    ```js
    {
      "client_id": "27f41c8da685c81a890f9e5f8ce48387",
      "client_secret": "78384ec8f5e9f098a18c586ad8c14f72",
      "redirect_uri": "https://3c6129e3.ngrok.io/oauth",
      "accounts_url": "https://accounts.livechatinc.com",
      "webhook_url": "https://3c6129e3.ngrok.io/webhook",
      "webhook_secret": "arbitrary_string"
    }
    ```

- Start the app

    ```bash
    go run .
    ```

### Install LiveChat app privately

1. Sign in to [Developer Console](https://developers.livechatinc.com/console/).
2. Go to the [Apps](https://developers.livechatinc.com/console/apps/) section.
3. Select the app created in one of the previous sections.
4. Move to the Private installation section and hit the install button.

**Attention!** For production applications, you want to publish the app.

### Post-install
Go to `https://direct.lc.chat/$YOUR_LICENSE_ID`, start a chat, say hi and wait for Echo to respond.
