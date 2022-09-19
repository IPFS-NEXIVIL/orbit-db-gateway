# go-orbit-db-gateway

## Server start

```
go run .
```

## Data insert to orbit db

`POST` `http://localhost:8001/ipfs/paste`

- example response

```
{
    "content": "test content",
    "project": "example project"
}
```

## Data get from orbit db

`GET` `http://localhost:8001/ipfs/get/{id}`

- example response

```
{
    "content": "test content",
    "date": "2022-08-24 15:34:40",
    "id": "d4c03a84-2376-11ed-b61b-3ccd36628b4c",
    "project": "example project"
}
```
