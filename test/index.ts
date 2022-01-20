import "dotenv/config"
import {RestClient} from "@keiryo/rest";
import consola from "consola";

import {
    mentionUserCommand,
    pingChatInputCommand,
    poopButton,
    repeatMessageInputCommand, youtubeAutocomplete,
    youtubeChatInputCommand
} from "./handlers";
import {createAmqp, createConsumer, ResponseOptions} from "./amqp";

import {
    APIApplicationCommandAutocompleteResponse,
    APIChatInputApplicationCommandInteraction,
    APIInteraction,
    APIInteractionResponse,
    APIMessageApplicationCommandInteraction,
    APIMessageComponentInteraction,
    APIUserApplicationCommandInteraction,
    ApplicationCommandOptionType,
    ApplicationCommandType,
    InteractionType,
    RESTPostAPIChatInputApplicationCommandsJSONBody,
} from "discord-api-types/v9";
import type {
    APIApplicationCommandAutocompleteInteraction
} from "discord-api-types/payloads/v9/_interactions/autocomplete";

const REST = new RestClient();
REST.token = process.env.DISCORD_TOKEN!!

const HANDLERS: InteractionHandlers = {
    [InteractionType.MessageComponent]: {
        "poop": poopButton
    },
    [InteractionType.ApplicationCommandAutocomplete]: {
        "youtube": youtubeAutocomplete
    },
    [InteractionType.ApplicationCommand]: {
        [ApplicationCommandType.User]: {"Mention": mentionUserCommand,},
        [ApplicationCommandType.ChatInput]: {
            "ping": pingChatInputCommand,
            "youtube": youtubeChatInputCommand
        },
        [ApplicationCommandType.Message]: {"Repeat": repeatMessageInputCommand}
    }
}

async function main() {
    consola.info("launching kantoku tester")

    if (process.argv.includes("--create")) {
        await REST.post(`/applications/${process.env.APPLICATION}/guilds/${process.env.GUILD}/commands`, {
            body: {
                type: ApplicationCommandType.ChatInput,
                name: "youtube",
                options: [
                    {
                        type: ApplicationCommandOptionType.String,
                        name: "query",
                        autocomplete: true,
                        description: "The query to search for.",
                        required: true
                    }
                ],
                description: "Searches for youtube videos.",
                default_permission: true
            } as RESTPostAPIChatInputApplicationCommandsJSONBody
        })
    }

    /* create amqp. */
    const amqp = await createAmqp("127.0.0.1");
    await amqp.channel.assertQueue("poop")
    await amqp.channel.bindQueue("poop", "gateway", "INTERACTION_CREATE");

    consola.info("connected to amqp")

    /* create consumer. */
    const consumer = createConsumer(amqp.channel, consumeInteraction)
    await amqp.channel.consume("poop", consumer);

    consola.info("now consuming interactions.")
}

async function consumeInteraction(data: APIInteraction, response: ResponseOptions) {
    consola.info(`consumed interaction:`, data.id)
    if (data.type === InteractionType.Ping) {
        return response.ack()
    }

    let resp: APIInteractionResponse | undefined
    switch (data.type) {
        case InteractionType.MessageComponent:
            resp = await HANDLERS[InteractionType.MessageComponent][data.data.custom_id]?.(data);
            break;
        case InteractionType.ApplicationCommand:
            resp = await HANDLERS[InteractionType.ApplicationCommand][data.data.type]?.[data.data.name]?.(data as any)
            break
        case InteractionType.ApplicationCommandAutocomplete:
            resp = await HANDLERS[InteractionType.ApplicationCommandAutocomplete][data.data!!.name]?.(data)
            break
    }

    if (resp) {
        await response.reply(resp)
    }

    await response.ack()
}

void main();

type InteractionHandlerRecord<I extends APIInteraction, R extends APIInteractionResponse = APIInteractionResponse> = Record<
    string,
    (interaction: I) => R | Promise<R | undefined> | undefined
>

interface InteractionHandlers {
    [InteractionType.MessageComponent]: InteractionHandlerRecord<APIMessageComponentInteraction>
    [InteractionType.ApplicationCommandAutocomplete]: InteractionHandlerRecord<APIApplicationCommandAutocompleteInteraction, APIApplicationCommandAutocompleteResponse>
    [InteractionType.ApplicationCommand]: {
        [ApplicationCommandType.User]: InteractionHandlerRecord<APIUserApplicationCommandInteraction>
        [ApplicationCommandType.Message]: InteractionHandlerRecord<APIMessageApplicationCommandInteraction>
        [ApplicationCommandType.ChatInput]: InteractionHandlerRecord<APIChatInputApplicationCommandInteraction>
    }
}
