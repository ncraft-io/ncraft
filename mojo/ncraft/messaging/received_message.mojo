

type ReceivedMessage: Message {
    ack_id: String @11
    subscription: String @12
    publish_time:  Timestamp @13
}