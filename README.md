# norman

Realtime distributed OLAP datastore, designed to answer OLAP queries with low latency written in Go

## Swagger

To generate the swagger definition run

```
swag init -d "./internal/commander"
swag init -d "./internal/broker"
```

## cURLs requests for testing

To create a new Schema

```
curl -X POST -H "Content-Type: application/json" -d @test/commander/create_schema.json http://localhost:8080/commander/v1/tenants/schemas
```
