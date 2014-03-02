# PowerRest

## About

Restful API for PowerDNS. It's written in Go (golang) and is extremely light on resources.

## Limitations
* PowerRest supports MySQL and PostgreSQL. Use either `mysql` or `postgres` in the config file.
* No authentication support. It should be used only in trusted environments or behind a reverse proxy that has authentication.

## Examples using Curl

### Domains

Response

The operations all return `200 OK` on success, except the create operation which returns the location of the created domain with the `201 Created` status code. The create operation also returns the newly created object as JSON in the body.

List domains

`curl "http://127.0.0.1/v1/domains"`

Create new domain

`curl -X POST --data-binary '{ "name": "example.com" }' "http://127.0.0.1/v1/domains"`

    HTTP/1.1 201 Created
    Location: /v1/domains/1

Update a domain

`curl -X POST --data-binary '{ "name": "example.org" }' "http://127.0.0.1/v1/domains/1"`

Delete a domain

`curl -X DELETE "http://127.0.0.1/v1/domains/1"`

### Records

List records

`curl "http://127.0.0.1/v1/records"`

Create new record

`curl -X POST --data-binary '{"domain_id":1,"name":"example.com","type":"A","content":"192.168.1.1","ttl":3600}' "http://127.0.0.1/v1/records"`

Update record

`curl -X POST --data-binary '{"domain_id":1,"name":"example.com","type":"A","content":"192.168.1.2","ttl":3600}' "http://127.0.0.1/v1/records/1"`

Delete record

`curl -X DELETE "http://127.0.0.1/v1/records/1"`