# Design

I designed this system in two discrete parts. The first part is the consumer written
in python and was designed to be ran in parallel (multiple instances of python at the same time).
I wrote this in python because it's lightweight and mostly I/O bound. The consumer parses the
raw bytes from the source (in this case udp_emitter.js), and turns it into JSON containing all
of the information from the raw message. This JSON is POSTed to the second part of the system,
the server.

I wrote the server in Golang to take advantage of the goroutines and speed when doing CPU bound
work. The server listens for POST requests to the `/add` endpoint, which then parses the JSON
and add it to the database.

Every write/read/persist/restore action on the database is atomic to prevent any possible race conditions.

The database stores each transaction, and all of the byte information associated with each message.
After the 20 seconds has elapsed since the transaction was first recorded, the server logs whether or not
the entire message was received. The entire database is persisted to disk on a configurable and periodic basis.
If the process was stopped or killed, there would only be a data loss of less than or equal to that period, as
the server restores all data on startup.


# Run

- To run the consumers, `./run.sh`
- To run the server, `go run server.go`
- To run the server tests, `go test`

# Test run

Start the server and run `test_harness.py`. This will send 10 randomly created messages to the server.

# TODO

- Write more tests! A full end to end integration test, more unit tests, etc.

