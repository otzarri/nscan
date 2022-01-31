# nscan

Simple network scanner. I developed `nscan` to learn Go programming and to have an existing codebase for improving my skills in the future. Don't expect regular commits here.

## Build

Build nscan and test the binary:

```
$ go build
$ ./nscan -host example.com
```

## Run

To run nscan without building the binary file:

```
$ go run main.go -host example.com
```

## Testing

Run tests as usual in Go:

```
$ go test -v ./...
```

To run tests with coverage:

```
$ go test -v -coverprofile=coverage.out ./...
$ go tool cover -html=coverage.out
```
