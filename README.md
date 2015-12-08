
# fbsender
A stream data sender to named pipe or tcp socket

### INSTALL

```
$ go get github.com/digsy/fbsender
```

### USAGE

```
$ $GOPATH/bin/fbsender send.conf
```

### CONFIGURATION
- name : sender name
- transaction_per_sec : number of transaction to be sent in one second
- source : data to be sent
- destination : path to filesystem or socket address

```
# exmaple
{
  "configurations":[
    {
      "name":"sender1",
      "transaction_per_sec":10,
      "source":"./data/input",
      "destination":"./data/output1"
    },
    {
      "name":"sender2",
      "transaction_per_sec":5,
      "source":"./data/input",
      "destination":"127.0.0.1:20000"
    }
  ]
}
```
