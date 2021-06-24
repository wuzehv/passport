fmt:
	gofmt -l -w ./

client:
	./run.sh client &
	./run.sh client -addr=127.0.0.1:8082 &

clean:
	rm -f passport