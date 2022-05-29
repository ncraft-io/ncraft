
type Subscription {
    name: String @1

    topic: String @5
    group: String @6

    pull: Bool @7

    auto_ack: Bool @9
    ack_timeout: Duration @10

    endpoint: PushEndpoint @15
}
