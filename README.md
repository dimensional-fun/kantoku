![Kantoku](https://user-images.githubusercontent.com/44017640/164708096-1d0585d4-5963-4011-bb36-763d413f1acc.png)

## 📦 Features

- super-duper fast
- publishes interactions to a NATS subject
- interaction testing route

## ⛓️ usage

### to start kantoku

1. fill in the `kantoku.toml` configuration file.
2. run the server thingy
3. go to your application in the [**discord developer portal**](https://discord.com/developers/applications) and set
   your interactions endpoint url to `https://<domain>/v1/interactions`

### implementing kantoku into your code base

Whenever Discord `POST`s an interaction to `/v1/interactions` Kantoku will request an interaction response on the
configured NATS subject.

_wip_

###### [Discord Server](https://discord.gg/8R4d8RydT4)

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

## 📜 contributors

- [@melike2d](https://github.com/melike2d)
- [@Topi314](https://github.com/Topi314)

---

[Dimensional Fun](https://dimensional.fun) &bull; Licensed under [**LGPL-2.1**](/LICENSE) 
