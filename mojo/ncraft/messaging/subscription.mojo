
type Subscription {
    name: String @1

    topic: String @5
    group: String @6

    pull: Bool @7

    auto_ack: Bool @9
    ack_timeout: Duration @10
    pull_max_waiting: Int64 @11 // default is 128
    pending_msg_limit: Int64 @12
    pending_bytes_limit: Int64 @13

    endpoint: PushEndpoint @15
}
