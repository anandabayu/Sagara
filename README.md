## Run

type this command to run:

```Shell
go get
go run main.go
```

## Test

type this command to test:

```Shell
cd tests/modeltests
go test -v
```

## API Documentation

Follow this link to see collections.
https://documenter.getpostman.com/view/13003337/UV5dfFYB

## Docker

before running on Docker config .env file first

```Shell
docker build .
docker-compose up -d
```

too rebuild image with docker-compose

```Shell
docker-compose down --remove-orphans --volumes
docker-compose up --build
```
