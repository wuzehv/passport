.PHONY: doc clean

fmt:
	gofmt -l -w ./

client:
	./run.sh client &
	./run.sh client -addr=127.0.0.1:8082 &

doc:
	swag init -o doc

clean:
	rm -f passport