
## Section: Manual Tests
## These tests exercise various aspects of the program's design from the command line
## and are intended for developer use
dev-race-test:
	ab -n 2000 -c 30 -T 'multipart/form-data' -p test/data/post-request-data http://localhost:8080/hash&
	ab -n 2000 -c 30 -m GET http://localhost:8080/1&
	ab -n 2000 -c 30 -m GET http://localhost:8080/stats&
	ab -n 2000 -c 30 -m GET http://localhost:8080/stats\?all&

dev-shutdown-stops-fielding-requests:
	curl —data “password=angryMonkey” http://localhost:8080/hash # Starts a job
	curl -X POST http://localhost:8080/shutdown
	curl http://localhost:8080/1 #Should be connection refused


## Section: Automated Tests
unit:
	go test
