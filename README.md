# 📡 mixtape-bot/kantoku

> http interactions server written in go

- publishes all interactions to rabbitmq
- probably bad because im very new to go

## ⛓️ usage

### to start kantoku

1. fill in the `kantoku.toml` configuration file.
2. run the server thingy
3. go to your application in the [**discord developer portal**](https://discord.com/developers/applications) and set
   your interactions endpoint url to `https://<domain>/v1/interactions`
4. it should look something like this when working!

![yessir](https://media.discordapp.net/attachments/830270945213284403/933854420410728458/unknown.png)

### implementing kantoku into your code base

make sure to look at my typescript implementation [**here**](https://github.com/mixtape-bot/kantoku-example).

> You **_WILL_** need to reply to the published RMQ message or else the interaction will fail,
> this is the equivalent of `POST`ing a response to the discord api.
>
> You **_WILL_** also need to acknowledge the message, or it'll be consumed multiple times.

###### [Discord Server](https://discord.gg/Vkbmb8kuH4)

## 📁 api

### `GET /v1`

#### Http Response
```json
{
    "data": "Hello, World!",
    "success": true
}
```

### `POST /v1/interactions`

#### Http Request

- `X-Signature-Ed25519` ed25519 signature
- `X-Signature-Timestamp` timestamp of the signature

<https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object>

#### Http Response

<https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object>

#### Amqp Response

- `body` is returned as the `Http Response`

```json
{
    "headers": {},
    "body": <binary data>
}
```

## 📜 credits

- [**disgo**](https://github.com/DisgoOrg/Disgo)
- [**spectacles**](https://github.com/spac-tacles/go)
- [**suggestionsbot/voting**](https://github.com/suggestionsbot/voting)

---

Mixtape Bot &bull; Licensed under [**LGPL-2.1**](/LICENSE) 
