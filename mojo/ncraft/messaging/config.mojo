

type Config {
    provider: String @1

    broker: String @2
    service_name: String @3

    subscriptions: [Subscription] @10

    nats: Nats @15
}

type Nats {
    jet_stream: String @1
    topic_names: [String] @2
    max_msgs: Int64 @3
    max_age: Int64 @4
}
