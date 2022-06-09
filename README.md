Fib server.

It is a hypothetical example of gRPC server to conduct some
longrunning operations. The server implements both FibService
and longrunning.OperationsService.  FibService receives a
request of calculating a fibonacci number, and then returns
the name of the operation immediately.

The operation runs in background and its status can be fetched
through OperationsService (right now it only implements
GetOperation for simplicity).

# How to build

```sh
% go build -o bin/server ./cmd/server
% go build -o bin/client ./cmd/client
```

# Sample usage

```sh
% ./bin/server &
% ./bin/client -n 45
```
