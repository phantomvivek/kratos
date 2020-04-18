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
The JSON test config (Please go through [sample config](https://github.com/phantomvivek/kratos/blob/master/config.json)) is divided into the following parts:

* config {}: Base config for app url & timeout value
  * url: The url of the websocket app (eg: ws://localhost:8080/)
  * timeout: Timeout, in seconds, before giving up on the connection (timeout errors will be reported)

* hitrate [ {} ]: An array of objects with each object containing the below keys:
  * duration: Duration for this step to be run, in seconds
  * end: The number of connections per second to end with at the end of duration.
  -- See example config for more info on how this works
  
* tests [ {} ]: An array of objects with each object containing a test for the connection made to the app
  * type: Can have the below possible values:
    * "message": To send a message
    * "sleep": Not do anything for a particular duration (only works with the duration argument)
    * "disconnect": To disconnect the socket connection to the app
  * send: When type is "message", the message to send. Can use variables from a CSV file (see dataFile variable).
          The index of the columns in CSV will be used as variables, like ${0}, ${1} & so on.
          For forming a message, only the same data row will be used, no two data rows will 
          contribute towards forming the same message.
  * replace: Boolean value, in case you don't want to replace constants in "send" string, 
             in case you want to use template variables in a message as is.
             
* dataFile: The path to the CSV file to use for data in the messages in tests.
            Do note that if you leave some columns in some rows blank then the replaced variable in "send"
            will be empty strings

* reporter: Kratos always reports to stdout. Use "type" as "statsd" to report to statsd as well.
  -- Below properties only work if type is "statsd"
  * type: Only supported value is "statsd". Will always report to stdout too.
  * host: Host for statsd daemon
  * port: Port for statsd daemon
  * prefix: Prefix string for all statsd metrics, eg: "myapp.loadtest"


## To Do:
- Support for yaml
- Support for custom reporter
- Context from app responses to be used in messages

## Inspiration:
[Vegeta](https://github.com/tsenart/vegeta) - A phenomenal load testing tool!

## License
See [LICENSE](LICENSE).
