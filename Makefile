INCLUDES+= -I=.
INCLUDES+= -I=$(GOPATH)/src/
INCLUDES+= -I=/usr/local/include

main: thyrc

thyrc: proto
	go build -o thyrc thyrc.go

proto: $(shell find . -name '*.proto' -type f | sed 's/proto$$/pb.go/')

./%.pb.go: %.proto
	protoc $(INCLUDES) --go_out=plugins=grpc,:. $(dir $<)*.proto

check: main
	go test $(shell find . -name '*_test.go' -type f | sort | xargs -n 1 dirname)

fmt:
	gofmt -w $(shell find . -name '*.go' -type f)

clean:
	rm -f thyrc
	rm -f $(shell find . -name '*.pb.go' -type f)

configure:
	./PREREQ.sh
