# Object Storage Reverse Proxy

This repository contains example of implementation [Kriechi/aws-s3-reverse-proxy](https://github.com/Kriechi/aws-s3-reverse-proxy) package with some additional feature like Virtual or Path style access convertion. This example is using a codebase from [dimaskiddo/codebase-go-cli](https://github.com/dimaskiddo/codebase-go-cli).

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.
See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Prequisites packages:
* Go (Go Programming Language)
* GoReleaser (Go Automated Binaries Build)
* Make (Automated Execution using Makefile)

Optional packages:
* Docker (Application Containerization)

### Deployment

#### **Using Container**

1) Install Docker CE based on the [manual documentation](https://docs.docker.com/desktop/)

2) Run the following command on your Terminal or PowerShell
```sh
docker run -d \
  -p 9000:9000 \
  -e OBJECT_STORAGE_PROXY_ENDPOINT="OBJECT_STORAGE_ENDPOINT" \
  -e OBJECT_STORAGE_PROXY_ACCESS_KEY="OBJECT_STORAGE_ACCESS_KEY" \
  -e OBJECT_STORAGE_PROXY_SECRET_KEY="OBJECT_STORAGE_SECRET_KEY" \
  -e OBJECT_STORAGE_PROXY_REGION="OBJECT_STORAGE_REGION" \
  --name object-storage-proxy \
  --rm dimaskiddo/object-storage-proxy:latest
```

#### **Using Pre-Build Binaries**

1) Download Pre-Build Binaries from the [release page](https://github.com/dimaskiddo/object-storage-proxy/releases)

2) Extract the zipped file

3) Run the pre-build binary
```sh
# MacOS / Linux
chmod 755 object-storage-proxy
./object-storage-proxy help
./object-storage-proxy proxy --help
./object-storage-proxy proxy --endpoint "OBJECT_STORAGE_ENDPOINT" --acccess-key "OBJECT_STORAGE_ACCESS_KEY" --secret-key "OBJECT_STORAGE_SECRET_KEY" --region "OBJECT_STORAGE_REGION"

# Windows
# You can double click it or using PowerShell
.\object-storage-proxy.exe help
.\object-storage-proxy.exe proxy --help
.\object-storage-proxy.exe proxy --endpoint "OBJECT_STORAGE_ENDPOINT" --acccess-key "OBJECT_STORAGE_ACCESS_KEY" --secret-key "OBJECT_STORAGE_SECRET_KEY" --region "OBJECT_STORAGE_REGION"
```

#### **Build From Source**

Below is the instructions to make this codebase running:
* Create a Go Workspace directory and export it as the extended GOPATH directory
```
cd <your_go_workspace_directory>
export GOPATH=$GOPATH:"`pwd`"
```
* Under the Go Workspace directory create a source directory
```
mkdir -p src/github.com/dimaskiddo/object-storage-proxy
```
* Move to the created directory and pull codebase
```
cd src/github.com/dimaskiddo/object-storage-proxy
git clone -b master https://github.com/dimaskiddo/object-storage-proxy.git .
```
* Run following command to renew and pull dependecies package
```
make vendor
```
* Until this step you already can run this code by using this command
```
make run
```

## Running The Tests

Currently the test is not ready yet :)

## Deployment

To build this code to binaries for distribution purposes you can run following command:
```
make release
```
The build result will shown in build directory
Or use Docker Images available in dimaskiddo/object-storage-proxy

## Built With

* [Go](https://golang.org/) - Go Programming Languange
* [GoReleaser](https://github.com/goreleaser/goreleaser) - Go Automated Binaries Build
* [Make](https://www.gnu.org/software/make/) - GNU Make Automated Execution
* [Docker](https://www.docker.com/) - Application Containerization

## Authors

* **Kriechi** - *Initial Source* - [Kriechi](https://github.com/Kriechi)
* **Dimas Restu Hidayanto** - *Refactorer* - [DimasKiddo](https://github.com/dimaskiddo)

See also the list of [contributors](https://github.com/dimaskiddo/object-storage-proxy/contributors) who participated in this project

## Annotation

You can seek more information for the make command parameters in the [Makefile](https://raw.githubusercontent.com/dimaskiddo/object-storage-proxy/master/Makefile)

## License

Copyright (C) 2023 Dimas Restu Hidayanto

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.