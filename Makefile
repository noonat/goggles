ifndef GOPATH
$(error GOPATH is not set)
endif

build: $(GOPATH)/bin/examples $(GOPATH)/bin/cube.js $(GOPATH)/bin/obj.js

clean:
	rm -f $(GOPATH)/bin/examples
	rm -f $(GOPATH)/bin/cube.js{,.map}
	rm -f $(GOPATH)/bin/obj.js{,.map}

$(GOPATH)/bin/examples: examples/*.go
	go install github.com/noonat/goggles/examples

$(GOPATH)/bin/cube.js: examples/cube/*.go *.go
	gopherjs install -m github.com/noonat/goggles/$(dir $<)

$(GOPATH)/bin/obj.js: examples/obj/*.go *.go
	gopherjs install -m github.com/noonat/goggles/$(dir $<)

run: build
	$(GOPATH)/bin/examples

watch: build
	gopherjs install -m -v -w \
		github.com/noonat/goggles/examples/cube \
		github.com/noonat/goggles/examples/obj

.PHONY: clean watch
