include makeconfig.env

.PHONY: clean init build debug env test

LDFlags := "-X '$(PACKAGE)/tokenguy/cmd.Version=$(VERSION)' $\
-X '$(PACKAGE)/tokenguy/cmd.CommitHash=$(COMMIT_HASH)' $\
-X '$(PACKAGE)/tokenguy/cmd.BuildTime=$(BUILD_TIMESTAMP)'"

build: bin/tokenguy

debug: DebugPrefix:=-dbg
debug: DebugFlags:=-gcflags=all="-N -l"
debug: bin/tokenguy$(DebugPrefix)

clean:
	rm bin/*

env:
	nix-shell

bin/tokenguy$(DebugPrefix): makeconfig.env go.mod main.go tokenguy/server.go tokenguy/cmd/root.go tokenguy/cmd/version.go tokenguy/cmd/server.go
	go get -d ./...
	go build -o bin/tokenguy$(DebugPrefix) -ldflags=$(LDFlags) $(DebugFlags) .
	chmod +x bin/tokenguy$(DebugPrefix)

docker:
	docker build -t $(REPOSITORY):$(VERSION) -t $(REPOSITORY):$(COMMIT_HASH) .
