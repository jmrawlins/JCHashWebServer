# JCHashWebServer
A password (or any string, really) hashing https server.

## Building
The program is written using the Go standard library alone in go 1.17, so there are no external go package dependencies.

Build into working directory using `go build`
Build into $GOPATH/bin using `go install`

## Testing
Unit test using `make test`

## Running
`./JCHashWebServer` for usage

## Endpoints
The server has the following endpoints:

- /hash
  - Request: POST using form data
  - Returns: hash retrieval id. Hashes can be retrieved using the `/<hash_id>` endpoint
  - Creates a hash job to hash the requested string in 5 seconds.
  - ex: curl --data "password=angryMonkey" http://localhost:8080/hash

- /<hash_id>
  - Request: GET
  - Returns: the hashed password associated with the requested hash_id

- /stats
  - Request: GET
  - Returns: request handling statistics for the /hash endpoint in json format

- /stats?all
  - Request: GET
  - Returns: request handling statistics for all endpoints since the server started, in json format

- /shutdown
  - Request: POST
  - Action: Signals the server to shut down gracefully
  - Returns: just 200 response code

## Assumptions
- Production-ready requirements:
My personal requirements for this minimum viable product were:
- Unit tests for packages
  - Good coverage with some error checking. Ideally for CICD we'd verify them all but I didn't have time.
- System tests for the program as a whole
  - They exercise all regular use cases while also stress-testing it a bit
- Logging
  - Very strong on errors, but optional at INFO level initially.
- Persistence
  - A production ready server would have a persistent data store, so I did implement that.

This program is currently meant to run on unix/linux systems only - the filesystem io and makefile
are not meant for windows file systems.

# Notes
- Not scale-ready:
For now I've assumed this is meant to run as a single instance on a single server.
In order to scale it up, I would want to use a centralized data storage (e.g., a database).
Doing so would allow us to horizontally scale server processes behind a load balancer.

- System-level tests:
There are no integration tests and the system level testing is primitive and manual, which I
would want to change in a true production situation -- however, they do prove the server can
handle a decent amount of traffic and would (did) find race conditions if (when) they existed.

- No peer review:
Since this was my first Go program I wanted terribly to get it peer reviewed by a more
experienced Go developer, but I didn't in the spirit of this being a personal evaluation.
I would never write this much code professionally without having other eyes on it! I still
fear that my C++ roots show in this project and I need more practice thinking in Go.

## Improvements
In future iterations I would:
- write a separate job handling goroutine to receive hashing requests
and process them after 5 seconds, in order to remove the sleep from the handler code and
allow those package tests to run without waiting for the code to sleep.
- revisit the shutdown mechanism to standardize it using context and cancel functions.
- fill out info level logging of requests instead of logging mostly errors
- talk about design with someone more familiar with Go
