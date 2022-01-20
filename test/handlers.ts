import {request} from "undici";
import type {
    APIApplicationCommandAutocompleteResponse,
    APIApplicationCommandInteractionDataStringOption,
    APIChatInputApplicationCommandInteraction,
    APIInteraction,
    APIInteractionResponse,
    APIMessageApplicationCommandInteraction,
    APIUserApplicationCommandInteraction
} from "discord-api-types/v9";
import {ButtonStyle, ComponentType, InteractionResponseType} from "discord-api-types/v9";
import type {
    APIApplicationCommandAutocompleteInteraction
} from "discord-api-types/payloads/v9/_interactions/autocomplete";

export async function youtubeAutocomplete(i: APIApplicationCommandAutocompleteInteraction): Promise<APIApplicationCommandAutocompleteResponse | undefined> {
    const option = (i.data?.options?.[0] as APIApplicationCommandInteractionDataStringOption)
    if (!option.focused || !option.value) {
        return Promise.resolve(undefined)
    }

    const results = await request(`https://suggestqueries-clients6.youtube.com/complete/search?q=${encodeURIComponent(option.value)}&h=en&client=firefox&gl=us`)
    const [input, suggestions] = await results.body.json() as [string, string[]]

    return {
        type: InteractionResponseType.ApplicationCommandAutocompleteResult,
        data: {
            choices: [
                {
                    value: input,
                    name: input
                },
                ...suggestions.map(s => ({name: s, value: s}))
            ]
        }
    }
}

export function youtubeChatInputCommand(i: APIChatInputApplicationCommandInteraction): APIInteractionResponse {
    const query = i.data.options?.[0] as APIApplicationCommandInteractionDataStringOption;
    console.log(query);

    return {
        type: InteractionResponseType.ChannelMessageWithSource,
        data: {
            content: query.value,
        }
    }
}

export function poopButton(_: APIInteraction): APIInteractionResponse {
    return {
        type: InteractionResponseType.ChannelMessageWithSource,
        data: {
            flags: 64,
            content: "poopie"
        }
    }
}

export function mentionUserCommand(inter: APIUserApplicationCommandInteraction): APIInteractionResponse {
    return {
        type: InteractionResponseType.ChannelMessageWithSource,
        data: {
            content: `<@${inter.data.resolved.users[inter.data.target_id]?.id}>`
        }
    }
}

export function pingChatInputCommand(_: APIInteraction): APIInteractionResponse {
    return {
        type: InteractionResponseType.ChannelMessageWithSource,
        data: {
            embeds: [
                {
                    description: "Hi, i don't have a ping cus im stateless uwu",
                    color: 0xB963A5
                }
            ],
            components: [
                {
                    type: ComponentType.ActionRow,
                    components: [
                        {
                            type: ComponentType.Button,
                            custom_id: "poop",
                            label: "Hello",
                            style: ButtonStyle.Primary,
                            disabled: false
                        }
                    ]
                }
            ]
        }
    }
}

export function repeatMessageInputCommand(interaction: APIMessageApplicationCommandInteraction): APIInteractionResponse {
    return {
        type: InteractionResponseType.ChannelMessageWithSource,
        data: {
            content: interaction.data.resolved.messages[interaction.data.target_id].content
        }
    }
}
