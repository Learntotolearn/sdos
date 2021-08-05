export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)"

clean:
	rm -f sdos