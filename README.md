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
    "content": "content"
}
```

## Data get from orbit db

`GET` `http://localhost:8001/ipfs/get/{id}`

- example response

```
{
    "content": "testtesttestrimi",
    "date": 1661317218848,
    "id": "a664660e-2369-11ed-aa3c-3ccd36628b4c"
}
```
