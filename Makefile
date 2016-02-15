APP_NAME = "coyote"
RELEASE_ROOT = "release"
RELEASE_COYOTE = "release/coyote"
NOW = $(shell date -u '+%s')

build:
	go build main.go

pack: build
	rm -rf $(RELEASE_COYOTE)
	mkdir -p $(RELEASE_COYOTE)
	mv main $(RELEASE_COYOTE)/$(APP_NAME)
	cp -r examples/custom $(RELEASE_COYOTE)/
	cd $(RELEASE_ROOT) && zip -r $(APP_NAME)-$(NOW).zip $(APP_NAME)

release: pack
