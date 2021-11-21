# CustomerAPI

Small microservice written with Go, Gin Web Framework, gorm and gocron. It connects to local postgresql (`docker-compose up`).

## Ping
[GET] /ping

## Create customer
[POST] /api/clients

Example payload:

```json
{
  "email": "hello@example.com",
  "title": "lunch",
  "content": "dumplings",
  "mailing_id": 1
}
```

## Get client by ID
[GET] /api/clients/:id

## Delete client by ID
[DELETE] /api/clients/:id

## Get all clients
[GET] /api/clients

## Mail clients
[POST] /api/clients/send

After logging a message, deletes all clients with the given `mailing_id`.
