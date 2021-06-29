.PHONY: doc clean

swag = swag init -o doc
fmt = gofmt -l -w ./

run:
	$(fmt)
	$(swag)
	go run main.go

fmt:
	$(fmt)

init:
	go run console/init.go

client1:
	go run console/client.go

client2:
	go run console/client.go -addr=127.0.0.1:8082

doc:
	$(swag)

clean:
	rm -f passport