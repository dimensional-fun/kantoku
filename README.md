# 📡 mixtape-bot/kantoku

> http interactions server written in go

- publishes all interactions to rabbitmq
- probably bad because im very new to go

## ⛓️ usage

### to start kantoku

1. fill in the `kantoku.toml` configuration file.
2. run the server thingy
3. go to your application in the [**discord developer portal**](https://discord.com/developers/applications) and set your interactions endpoint url to `https://<domain>/v1/interactions`
4. it should look something like this when working!

![yessir](https://media.discordapp.net/attachments/830270945213284403/933854420410728458/unknown.png)

### implementing kantoku into your code base

make sure to look at my typescript implementation [**here**](/test).

##### ! _note:_

You WILL need to reply to the published RMQ message or else the interaction will fail,
this is the equivalent of `POST`ing a response to the discord api.

You will need to also acknowledge the message, or it'll be consumed multiple times.

## 📜 credits

- [**disgo**](https://github.com/DisgoOrg/Disgo)
- [**spectacles**](https://github.com/spac-tacles/go)
- [**suggestionsbot/voting**](https://github.com/suggestionsbot/voting) 

---

Mixtape Bot &bull; Licensed under [**LGPL-2.1**](/LICENSE) 
