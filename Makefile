include makeconfig.env

.PHONY: clean init build debug env test deps-upgrade

LDFlags := "-X '$(PACKAGE)/tokenguy/cmd.Version=$(VERSION)' $\
-X '$(PACKAGE)/tokenguy/cmd.CommitHash=$(COMMIT_HASH)' $\
-X '$(PACKAGE)/tokenguy/cmd.BuildTime=$(BUILD_TIMESTAMP)'"

build: bin/$(ARCH)/tokenguy

debug: DebugPrefix:=-dbg
debug: DebugFlags:=-gcflags=all="-N -l"
debug: bin/$(ARCH)/tokenguy$(DebugPrefix)

clean:
	-rm -r bin/*

bin/%/tokenguy$(DebugPrefix): makeconfig.env go.mod main.go tokenguy/server.go tokenguy/validate.go tokenguy/cmd/root.go tokenguy/cmd/version.go tokenguy/cmd/server.go tokenguy/cmd/validate.go
	go get ./...
	GOARCH=$(ARCH) go build -o bin/$(ARCH)/tokenguy$(DebugPrefix) -ldflags=$(LDFlags) $(DebugFlags) .
	chmod +x bin/$(ARCH)/tokenguy$(DebugPrefix)

docker:
	docker build -t $(REPOSITORY):$(VERSION) -t $(REPOSITORY):$(COMMIT_HASH) .

deps-upgrade:
	go get -u
	go mod tidy
	echo "Don't forget to update golang in go.mod, .github/workflows/go.yml, and Dockerfile"

tag-commit:
	echo "git tag -a <version> -m <commit-msg>"