# authbaton

[![Build Status](https://github.com/xmidt-org/authbaton/actions/workflows/ci.yml/badge.svg)](https://github.com/xmidt-org/authbaton/actions/workflows/ci.yml)
[![codecov.io](http://codecov.io/github/xmidt-org/authbaton/coverage.svg?branch=main)](http://codecov.io/github/xmidt-org/authbaton?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmidt-org/authbaton)](https://goreportcard.com/report/github.com/xmidt-org/authbaton)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=xmidt-org_authbaton&metric=alert_status)](https://sonarcloud.io/dashboard?id=xmidt-org_authbaton)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/xmidt-org/authbaton/blob/main/LICENSE)
[![GitHub Release](https://img.shields.io/github/release/xmidt-org/authbaton.svg)](CHANGELOG.md)

## Summary
AuthBaton is an authentication service for applications behind a reverse proxy.
## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Details](#details)
- [Usage](#usage)
- [Build](#build)
- [Deploy](#deploy)
- [Contributing](#contributing)

## Code of Conduct

This project and everyone participating in it are governed by the [XMiDT Code Of Conduct](https://xmidt.io/docs/community/code_of_conduct/). 
By participating, you agree to this Code.

## Details
AuthBaton is meant to be used as a helper authentication microservice to reverse proxy tools such as NGINX.

The diagram below shows the path that a request follows before reaching the protected application.  
![Diagram](docs/diagrams/Auth-baton%20Success%20Auth%20Flow.png)
## Usage
```
curl http://localhost:6800 -i
HTTP/1.1 403 Forbidden
X-Server-Name: authbaton
X-Server-Version: development
Date: Mon, 05 Apr 2021 21:18:24 GMT
Content-Length: 0
Connection: close
```

```
curl http://localhost:6800/original/request/path -H "Authorization: Basic dXNlcjpwYXNz" -i
HTTP/1.1 200 OK
X-Server-Name: authbaton
X-Server-Version: development
Date: Mon, 05 Apr 2021 21:21:46 GMT
Content-Length: 0
Connection: close
```
**Note:** AuthBaton accepts any URL path. This allows bascule capability checks 
to work properly as the reverse proxy can simply reuse the URL path of the original request.

## Build
### Source
In order to build from source, you need a working 1.x Go environment. Find more information on [Go website](https://golang.org/doc/install).

Then, clone the repo and build using make:

```bash
git clone git@github.com:xmidt-org/authbaton.git
cd authbaton
make build
```

### Makefile

The Makefile has the following options you may find helpful:
* `make build`: builds the authbaton binary
* `make test`: runs unit tests with coverage for authbaton
* `make clean`: deletes previously-built binaries and object files

### RPM

First have a local clone of the source and go into the root directory of the 
repository.  Then use rpkg to build the rpm:
```bash
rpkg srpm --spec <repo location>/<spec file location in repo>
rpkg -C <repo location>/.config/rpkg.conf sources --outdir <repo location>'
```

## Deploy
Once the binary is built, run:
```
./authbaton
```
Ensure that the `authbaton.yaml` config file is in one of the following folders:
- The current working directory
- `$HOME/.authbaton`
- `/etc/authbaton`


### Supported Reverse Proxies
We currently have an example configuration file only for NGINX. However, any reverse proxy that can authenticate an external request by consulting authbaton is supported.

See example configurations [here](/docs/example-config). We are happy to take contributions for example config files for other reverse proxies. 

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md).
