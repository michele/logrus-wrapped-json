## WrappedJSONFormatter

This [logrus](https://github.com/sirupsen/logrus) formatter allows you to wrap you JSON logs into a first level element named after the field `kind`.

An example with the standard `JSONFormatter` would be:

```json
{"animal":"walrus","level":"info","msg":"A giant walrus appears!",
"size":10,"time":"2014-03-10 19:57:38.562500591 -0400 EDT"}
```

With `WrappedJSONFormatter`, you could pass an extra field like this:

```go
log = log.WithField("kind", "wrapped_key_name")
```

And the output would become:

```json
{"wrapped_key_name": {"animal":"walrus","level":"info","msg":"A giant walrus appears!",
"size":10,"time":"2014-03-10 19:57:38.562500591 -0400 EDT"}}
```

Example:

```go
package main

import (
	log "github.com/sirupsen/logrus"
	wrapped "github.com/michele/logrus-wrapped-json"
)

func init() {
  log.SetFormatter(&wrapped.WrappedJSONFormatter{})
}

func main() {
    httpLog := log.WithField("kind", "response_log")
    dbLog := log.WithField("kind", "sql_log")

    ...
}
```
