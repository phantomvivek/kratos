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

This lib has evolved from being a measly node.js script to a somewhat good (hopefully) tool to test websocket applications. Suggestions / PRs / Questions welcome!

---

## API:
The JSON test config ([sample config](https://github.com/phantomvivek/kratos/blob/master/config.json)) is divided into the following parts:

* config: Base config for app url & timeout value
  * url: The url of the websocket app (eg: ws://localhost:8080/)
  * timeout: Timeout, in seconds, before giving up on the connection (timeout errors will be reported)


* hitrate: An array of objects with each object containing the below keys:
  * duration: Duration for this step to be run, in seconds
  * end: The number of connections per second to end with at the end of duration.
  -- See example config for more info on how this works


* tests: An array of objects with each object containing a test for the connection made to the app
  * type: Can have the below possible values:
    * "message": To send a message
    * "sleep": Not do anything for a particular duration (only works with the duration argument)
    * "disconnect": To disconnect the socket connection to the app
  * send: When type is "message", the message to send. Can use variables from a CSV file (see dataFile variable). The index of the columns in CSV will be used as variables, like ${0}, ${1} & so on. For forming a message, only the same data row will be used, no two data rows will contribute towards forming the same message.
  * replace: Boolean value, in case you don't want to replace constants in "send" string, in case you want to use template variables in a message as is.


* dataFile: The path to the CSV file to use for data in the messages in tests. A connection will use data from only a single row for its `tests`


* reporter: Kratos always reports to stdout. Use "type" as "statsd" to report to statsd as well. (Below properties only work if type is "statsd")
  * type: Only supported value is "statsd". Will always report to stdout too.
  * host: Host for statsd daemon
  * port: Port for statsd daemon
  * prefix: Prefix string for all statsd metrics, eg: "example.myapp"

  Sample reporter example JSON for statsd:
  ```json
  "reporter": {
    "type": "statsd",
    "host": "localhost",
    "port": 8125,
    "prefix": "example.myapp"
  }
  ```

### API Example
Consider the following example for how hitrate & tests work. First, we will look at the hitrate array:
```javascript
hitrate: [
  { duration: 10, end: 20 },
  { duration: 10, end: 20 },
  { duration: 10, end: 0 }
]
```
We will assume the app isn't closing any connections made to it:
1. For the first phase of 10 seconds, start with making 2 connections in the first second, increase 2 connections *per second*, and end with 20 connections *per second*. This is the *ramp-up* duration.
    - 1st second = Make 2 connections = App has 2 active connections
    - 2nd second = Make 4 connections = App now has 6 active connections
    ...
    - 10th second = Make 20 connections = App now has 110 active connections
    - Use the arithmetic sequence sum formula = `n/2 * (2a + (n - 1)d)`, where *n = number of terms (10 in this example)*, *a = first term (2, connections in the 1st second)* & *d = difference (2, per second increment)*
2. For the next phase of 10 seconds, keep making 20 connections *per second* as `end=20`. This is the *steady* duration. In this duration, the app will receive a total of 200 connections (20 * 10).
3. For the 3rd phase of 10 seconds, decrease the connections made each second to reach 0 connections at the end of this duration. This is the *ramp-down* duration.

**Note:** This is just an illustrative example. Feel free to use any number of phases, durations or connections.


Next up, the tests array:
```javascript
tests: [{
  "type": "message",
  "send": "CSV data! ${0} and ${1} and ${2}",
  "replace": true
},
{
  "type": "message",
  "send": "Strings in this message won't be replaced! ${1} and ${2}"
},
{
  "type": "sleep",
  "duration": 2
},
{
  "type": "disconnect"
}]
```
For ***each connection*** made to the app, the connection will do the following *tests* in the order they're defined:
1. Send a message with variables replaced by CSV data
2. Send a message where the string is sent as is
3. Do nothing for 2 seconds
4. Disconnect the socket. (*Note: If any tests are defined after disconnecting the socket, it will error out*)

Example logs from the app being tested:
```
Socket opened
Message on socket: "CSV data! row1A and  row1B and  row1C"  -- Note that the variables were replaced from CSV
Message on socket: "Strings in this message won't be replaced! ${1} and ${2}"
...nothing for 2 seconds...
Socket disconnected

Socket opened
Message on socket: "CSV data! row2A and  row2B and  row2C"   -- Note that the variables were replaced from CSV
Message on socket: "Strings in this message won't be replaced! ${1} and ${2}"
...nothing for 2 seconds...
Socket disconnected
```

---

## To Do:
- Support for yaml
- Support for custom reporter
- Context from app responses to be used in messages

---

## License
See [LICENSE](LICENSE).
