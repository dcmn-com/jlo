# jlo

Light-weight json logging in go

![ ](https://media.giphy.com/media/l41lI61owQXsxHnvG/giphy.gif)

- default log level: `INFO`

## Usage

```go

l := jlo.NewLogger(os.Stdout)
l.Infof("I'm real")

l.SetLogLevel(jlo.DebugLevel)
l.Debugf("what you get is what you %s", "see")

l = l.WithField("@request_id", "aa33ee55")
l.Errorf("What you tryna to do to me?")

```

## Example output

```json
{"@timestamp": "2018-08-17T13:24:08.856339554Z","@level":"info","@message": "I'm real"}
{"@timestamp": "2018-08-17T13:24:08.856339554Z","@level":"debug","@message": "what you get is what you see"}
{"@timestamp": "2018-08-17T13:24:08.856391733Z","@level":"error","@request_id":"aa33ee55","@message": "What you tryna to do to me?"}
...
```
