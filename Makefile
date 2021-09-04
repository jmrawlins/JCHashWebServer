
dev-race-test:
	ab -n 8000 -c 30 -T 'multipart/form-data' -p test/data/post-request-data http://127.0.0.1:8080/hash&
	ab -n 8000 -c 30 -m GET http://127.0.0.1:8080/1&