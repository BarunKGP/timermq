# timer-mq: A message queue with timers, written in Golang

Message queues are systems that transport messages from a source to a receiver.
It is meant to achieve fault tolerance in highly distributed systems, through decoupling.
`TimerMQ` is a Golang implementation of a message queue that supports adding a timeout before the message is sent.
When the timeout expires, the message is sent to the receiver.
The message may also be cancelled any time before the timeout expires.

## How it works

1. Source sends a message to `TimerMQ` with a specified `delayMs` (0s by default, i.e., this behaves like a typical message queue)
2. `TimerMQ` returns a confirmation message with `messageId` for the newly created message.
3. If you cancel a message before its timeout expires, it is sent to a dead-letter queue and never makes its way to the receiver
4. When the timeout of a message expires (and if it hasn't been cancelled prior to that), the message is sent to the receiver.

## Supported Commands

-. `PUSH <value> <args>`: Pushes a string `value` into `TimerMQ` with default delay of `0ms`. This method will return the message id `id` of the the pushed message
-. `GET <id>`: Retrives a message with `id`. Returns an empty value if it does not exist.
-. `CANCEL <id>`: Cancels the message with id `id` if it is scheduled to be published and has not expired yet.

### Optional Args

When pushing a message into `TimerMQ`, you can set optional args modifying the behavior of the message

| Optional Arg | Supported Commands | Description                                                                                                                                  | Default |
| ------------ | ------------------ | -------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `delayMs`    | `PUSH`             | Sets a delay in milliseconds after which message will be pushed to `TimerMQ`. This message will be stored immediately if `durable` is `true` | 0 ms    |
| `durable`    | `PUSH`             | If `true`, the message will be stored in the persistence layer (in-memory, database, file, etc.)                                             | `false` |

## Messages

Messages are a composite, atomic unit of transaction within `TimerMQ`.
A message consists of a `<COMMAND>` and `val` specifying the action to be taken on `val`.
This behavior may be modified by passing along optional `args`.
Each message is separated by the newline characted (`\n`)

Each message sent to `TimerMQ` must be of the form:

```
<COMMAND> <val> <args>
```

`val` is then processed according to the specified `COMMAND` and `args`.

Once a message is passed, it cannot be modified.

### TODO:

- [ ] If persistence is enabled, each message is stored in the specified persistence layer.
- [ ] If logging is enabled, each message is logged to the specified log stream.

As of today, TimerMQ operates using a simple message protocol over TCP.
Support for other protocols such as MQTT and AMQP are in development.
