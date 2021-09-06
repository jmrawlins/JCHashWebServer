
## Section: Manual Tests
## These tests exercise various aspects of the program's design from the command line
## and are intended for developer use
dev-race-test:
	ab -l -v 2 -n 4000 -c 30 -T 'application/x-www-form-urlencoded' -p test/data/post-request-data http://localhost:8080/hash 2>&1 > /tmp/hashReport.log &
	ab -l -v 2 -n 4000 -c 30 -m GET http://localhost:8080/1 2>&1 > /tmp/getHashReport.log &
	ab -l -v 2 -n 4000 -c 30 -m GET http://localhost:8080/stats 2>&1 >  /tmp/statsReport.log &
	ab -l -v 2 -n 4000 -c 30 -m GET http://localhost:8080/stats\?all 2>&1 > /tmp/statsAllReport.log &

dev-shutdown-stops-fielding-requests:
	curl —data “password=angryMonkey” http://localhost:8080/hash # Starts a job
	curl -X POST http://localhost:8080/shutdown
	curl http://localhost:8080/1 #Should be connection refused


## Section: Automated Tests
unit:
	go test -race $$(go list)/datastore $$(go list)/handlers
