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
