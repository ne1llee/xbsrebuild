NAME=xbsrebuild
BINDIR=bin
GOBUILD=go build -ldflags '-w -s'

PLATFORM_LIST = \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64 \
	linux-arm64
WINDOWS_ARCH_LIST = \
	windows-386 \
	windows-amd64

all: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

windows-386:
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe
	
windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe


gz_releases=$(addsuffix .gz, $(PLATFORM_LIST))
zip_releases=$(addsuffix .zip, $(WINDOWS_ARCH_LIST))

$(gz_releases): %.gz : %
	chmod +x $(BINDIR)/$(NAME)-$(basename $@)
	gzip -f -S .gz $(BINDIR)/$(NAME)-$(basename $@)

$(zip_releases): %.zip : %
	zip -m -j $(BINDIR)/$(NAME)-$(basename $@).zip $(BINDIR)/$(NAME)-$(basename $@).exe

all-arch: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

releases: $(gz_releases) $(zip_releases)


clean:
	rm $(BINDIR)/xbs*