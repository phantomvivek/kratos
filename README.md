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
The JSON test config ([sample config](https://github.com/phantomvivek/kratos/blob/master/config.json)) is divided into the following parts:

* config: Base config for app url & timeout value
  * url: The url of the websocket app (eg: ws://localhost:8080/)
  * timeout: Timeout **in seconds** before giving up on the connection (timeout errors will be reported)

* hitrate: An array of objects with each object containing the below keys:
  * duration: Duration for this step to be run, in seconds
  * end: The number of connections per second to end with at the end of duration.
  ```
  Eg: An object like {duration: 10, end: 50} would start with 0 connections 
  (since this is the beginning of the test), and each second increases the 
  number of connections per second made to the app by 5. So the flow will be like:
  - 1st second = 5 connections made in this second
  - 2nd second = 10 connections made in this second
  ...
  - 9th second = 45 connections made in this second
  - 10th second = 50 connections made in this second
  So the app will have received 275 connections at the 10th second 
  (Summation of arithmatic sequence formula is S = (n/2) × (2a + (n−1)d) where n = 10, a = 5 & d = 5)
  ```
  The test begins with 0 connections to the app & ramps up according to the hitrate config with each stage either ramping up or ramping down the connections. The sample config is an example for this.

## To Do:
- Support for yaml
- Support for custom reporter
- Context from app responses to be used in messages

## Inspiration:
[Vegeta](https://github.com/tsenart/vegeta) - A phenomenal load testing tool!

## License
See [LICENSE](LICENSE).
