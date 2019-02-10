# Documentation

## Run the App

Easiest is to use `docker-compose`:

```sh
docker-compose up
```

And the app will be running at the address `0.0.0.0:8080`.

Alternatively you can run it manually if you have Redis running on localhost:

```sh
GO111MODULE=on go run main.go
```

## API

### Create a new voucher

Let's create couple of vouchers:

```sh
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"foo","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"bar","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"qux","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"foo1","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"bar1","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"qux1","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"foo2","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"bar2","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"qux2","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"foo3","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"bar3","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
curl -XPOST 0.0.0.0:8080/vouchers -d '{"id":"qux3","name":"Save £20 at Tesco","brand": "Tesco","value": "20"}'
```

### List vouchers

```sh
curl 0.0.0.0:8080/vouchers
```

Which will return list of vouchers:

```json
{
	"offset": 0,
	"nextOffset": 0,
	"vouchers": [{
		"id": "qux",
		"name": "Save £20 at Tesco",
		"brand": "Tesco",
		"value": "20"
	}, {
		"id": "bar",
		"name": "Save £20 at Tesco",
		"brand": "Tesco",
		"value": "20"
	}, {
		"id": "foo",
		"name": "Save £20 at Tesco",
		"brand": "Tesco",
		"value": "20"
	}]
}
```

You can use `nextOffset` for pagination. If it is -1, it means there is no more results. Just append `?offset=$nextOffset` to go to the next page. You can also use `count` to specify number of results per page.

### Retrieve voucher by ID

```sh
curl 0.0.0.0:8080/vouchers/foo
```

### Delete voucher 

```sh
curl -XDELETE 0.0.0.0:8080/vouchers/foo
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

## Dependency Management

Since Go 1.11, a new recommended dependency management system is via [modules](https://github.com/golang/go/wiki/Modules).

This is one of slight weaknesses of Go as dependency management is not a solved problem. Previously Go was officially recommending to use the [dep tool](https://github.com/golang/dep) but that has been abandoned now in favor of modules.
