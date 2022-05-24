
type Subscription {
    name: String @1

    topic: String @5
    group: String @6

    auto_ack: Bool @9

    endpoint: PushEndpoint @10
}
