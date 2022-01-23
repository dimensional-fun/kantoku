const {connect} = require("amqplib");
const {decode, encode} = require("@msgpack/msgpack");
const {amqpConsume, amqpConnect} = require("./amqp");

const queue = "interactions";

async function main([, , ...commands] = process.argv) {
    /* connect to amqp. */
    const { conn, chan } = await amqpConnect()

    /* create the consumer queue. */
    await chan.assertQueue(queue);
    await chan.bindQueue(queue, "gateway", "INTERACTION_CREATE");

    /* start consuming. */
    for (let i = 0; i < commands.length; i++) {
        await startConsumer(chan, i + 1, commands[i]);
    }
}

async function startConsumer(chan, consumer, command) {
    console.log("Hi from consumer", consumer)

    /* start consuming */
    await amqpConsume(chan, queue, async (data, responder) => {
        if (data.type === 2) {
            if (data.data.name === command) {
                console.log(consumer, "-", data.data.name, "yes")

                await responder.reply({
                    type: 4,
                    data: {
                        content: "hi from consumer " + consumer
                    }
                })

                await responder.ack();

                return;
            }

            console.log(consumer, "-", data.data.name, "no")
        }

        await responder.nack({requeue: true});
    });
}

void main();
