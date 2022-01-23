const {connect} = require("amqplib");

/**
 * @param {unknown} data
 * @returns {Buffer}
 */
function msgpackEncode(data) {
    const encoded = require("@msgpack/msgpack").encode(data);
    return Buffer.from(encoded, encoded.byteOffset, encoded.byteLength);
}

/**
 * @param {Buffer} buffer
 * @returns {unknown}
 */
function msgpackDecode(buffer) {
    return require("@msgpack/msgpack").decode(buffer);
}

/**
 * @param {import("amqplib").Channel} channel
 * @param {import("amqplib").Message} msg
 * @returns {Responder}
 */
function createResponder(channel, msg) {
    return {
        ack: () => channel.ack(msg),
        nack: ({ requeue = false, aut = false }) => channel.nack(msg, aut, requeue),
        reject: (requeue = false) => channel.reject(msg, requeue),
        reply: data => channel.sendToQueue(msg.properties.replyTo, msgpackEncode(data), { correlationId: msg.properties.correlationId }),
    }
}

/**
 * @returns {Promise<{conn: import("amqplib").Connection, chan: import("amqplib").Channel}>}
 */
async function amqpConnect() {
    /* connect to amqp. */
    const conn = await connect("amqp://127.0.0.1")
        , chan = await conn.createChannel();

    return { conn, chan }
}

/**
 * @param {import("amqplib").Channel} channel
 * @param {string} queue
 * @param {function(data: unknown, responder: Responder): any} fn
 * @returns {Promise<void>}
 */
async function amqpConsume(channel, queue, fn) {
    await channel.consume(queue, msg => {
        if (!msg) return
        console.log(msg.properties)
        fn(msgpackDecode(msg.content), createResponder(channel, msg))
    })
}

module.exports = {
    amqpConnect,
    amqpConsume
}

/**
 * @typedef {{reject: () => Promise<void>, ack: () => Promise<void>, nack: ({requeue?: boolean, aut?: boolean}) => Promise<void>, reply: (data: any) => Promise<void>}} Responder
 */