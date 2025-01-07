## Start

```shell
go run .
```

## Test

Mocked method

```shell
curl -ls http://localhost:8080/api/v1/users/profile
```

Implemented method

```shell
curl -ls http://localhost:8080/api/v1/users/profile2
```

Empty response

```shell
curl -ls http://localhost:8080/api/v1/status
```
