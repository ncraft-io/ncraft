| field | type | format | required | default | description |
|---|---|---|---|---|---|
| `ackTimeout` | `string` | `Duration` | N |  | A Duration represents a signed, fixed-length span of time represented as a count of seconds andfractions of seconds at nanosecond resolution. It is independent of any calendar and concepts like "day" or "month".It is related to Timestamp in that the difference between two Timestamp values is a Duration andit can be added or subtracted from a Timestamp. Range is approximately +-10,000 years. |
| `autoAck` | `boolean` |  | N |  |
| `endpoint` | `ncraft.messaging.PushEndpoint` |  | N |  |  |
| `group` | `string` |  | N |  |
| `name` | `string` |  | N |  |
| `pendingBytesLimit` | `integer` | `Int64` | N |  |
| `pendingMsgLimit` | `integer` | `Int64` | N |  |
| `pull` | `boolean` |  | N |  |
| `pullMaxWaiting` | `integer` | `Int64` | N |  |
| `topic` | `string` |  | N |  |
