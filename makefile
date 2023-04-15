NAME=xbsrebuild
BINDIR=bin
GOBUILD=go build -ldflags '-w -s'

PLATFORM_LIST = \
	linux-amd64 \
	linux-arm64 \
	windows-386 \
	windows-amd64
	
all: $(PLATFORM_LIST)

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

windows-386:
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@
	
windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

all-arch: $(PLATFORM_LIST)

clean:
	rm $(BINDIR)/*