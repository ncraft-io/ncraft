
type Subscription {
    name: String @1

    topic: String @5
    auto_ack: Bool @6

    endpoint: PushEndpoint @10
}
