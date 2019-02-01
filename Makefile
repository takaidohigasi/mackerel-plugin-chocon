NAME="mackerel-plugin-chocon"
VERSION=0.0.1
.PHONY: all build dep clean

all: dep build 

build:
	go build -ldflags "-X main.Version=${VERSION}" -o ${NAME}

dep:
	dep ensure -v

fmt:
	go fmt ./...

dist:
	git archive --format tgz HEAD -o $(NAME)-$(VERSION).tar.gz --prefix $(NAME)-$(VERSION)/

update:
	dep ensure -update

clean:
	@rm ${NAME}

tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin master