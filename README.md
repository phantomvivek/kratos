# kratos
Load Testing Tool for Websockets, written in Go.

### Why?
Websocket servers exist in a lot of spaces, especially gaming. Load testing is an important testing facet for these applications. Over the last few years, after having thrown the kitchen sink at some apps trying to break them, there seemed to be no good way of getting this done.

<ins>The problems faced primarily included:</ins>
1. Message exchanges between tool & app
2. Testing more than just the connection holding capacity of the app
3. Strenous load tests without the lib crashing (:man_facepalming:)

<ins>The solutions tried over the years before building this:</ins>
- Artillery.io
- Jmeter
- Gatling

This lib has evolved from being a measly node.js script to a somewhat good (hopefully) tool to test websocket applications. Suggestions / PRs welcome!

## API:
<TODO>

## To Do:
- Support for yaml
- Support for custom reporter
- Context from app responses to be used in messages

## License
See [LICENSE](LICENSE).
