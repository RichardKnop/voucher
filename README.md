# Documentation

## Run the App

If you have local development environmet set up, you can just do:

```sh
GO111MODULE=on go run main.go
```

And the app will be running at the address `0.0.0.0:8080`.

Alternatively, you can use `docker-compose`:

```sh
docker-compose up
```

Or build an image with `docker`:

```sh
docker build -t voucher .
```

And run it:

```sh
docker run --publish 8080:8080 --name test --rm voucher

```

## API

### Create a new voucher

curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"foo","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'

### Retrieve voucher by ID

curl 0.0.0.0:8080/vouchers/foo

## Development Setup

Make sure you have Golang development environment configured properly on your computer. You can follow the [Getting Started](https://golang.org/doc/install) guide.

Then you can pull the code from GitHub:

```sh
go get -u github.com/RichardKnop/voucher
```

You might also need `docker` and `protobuf`. On Macs, you can install them using homebrew:

```sh
brew install protobuf

brew cask install docker
ln -s /Applications/Docker.app/Contents/Resources/bin/* /usr/local/bin/
```

Also install Go wrappers for protobuf:

```sh
go get -u -v github.com/golang/protobuf/proto
go get -u -v github.com/golang/protobuf/protoc-gen-go
```

If you make any changes to the `pb/pb.proto` file, you can regenerate Golang structs like this:

```sh
protoc --go_out=paths=source_relative:. pb/*.proto
```

## Dependency Management

Since Go 1.11, a new recommended dependency management system is via [modules](https://github.com/golang/go/wiki/Modules).

This is one of slight weaknesses of Go as dependency management is not a solved problem. Previously Go was officially recommending to use the [dep tool](https://github.com/golang/dep) but that has been abandoned now in favor of modules.
