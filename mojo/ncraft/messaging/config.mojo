

type Config {
    provider: String @1

    broker: String @2
    service_name: String @3

    subscriptions: [Subscription] @10
}
