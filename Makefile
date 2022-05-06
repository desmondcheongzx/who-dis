BINS = who-dis

.PHONY: all $(BINS) clean test format

all: $(BINS)

$(BINS):
	@[ -f ./go.mod ] || go mod init github.com/desmondcheongzx/who-dis
	@go mod tidy
	go build ./cmd/$@

clean:
	rm -f $(BINS) *.lnx

test:
	go test ./test/* -v -race

format:
	gofmt -s -w .
