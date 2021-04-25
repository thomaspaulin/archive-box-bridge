# Request to Archive Box Bridge
A Go server which accepts URLs and passes them onto Archive Box for archiving.

## Requirements
- `docker-compose` is installed and present on the PATH.
- Archive Box is installed using docker-compose.
- The `ARCHIVE_BOX_DIR` environment variable is set to the ArchiveBox directory. That is, the directory you set up ArchiveBox in and which contains a `docker-compose.yml` file.
- Go is set up to build the binary.

## Running
1. `go build main.go`
2. Set `ARCHIVE_BOX_DIR` to the directory mentioned above
   1. Optionally set `ARCHIVE_BRIDGE_PORT` if you wish to run the server on a port other than 3344
3. Run the produced binary