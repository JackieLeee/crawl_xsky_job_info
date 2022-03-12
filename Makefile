export tag=v1.0

build:
	echo "building crawl_xsky_job_info binary"
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/crawl_xsky_job_info_linux
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/crawl_xsky_job_info_macos
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/crawl_xsky_job_info_windows.exe