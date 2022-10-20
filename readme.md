# # Simple REST API for Note-taking App witg auth 

### To run app

Setup your local app.env @ ./configs/app.env file and then run:

```bash
make run
```

or

```bash
docker compose --env-file .\configs\app.env up --build post-app
```
To run with gRPC audit logger please run with grpc server @ https://github.com/Arkosh744/grpc-audit-log
_________________________________________________
# TESTs

To run tests:

```bash
make test
#or
go test -v -count=1 ./...
```



# REST API:

### Examples of usage

#### First we need to sign-up:

`POST /auth/sign-up`

```json
{
  "name": "NewGopher",
  "email": "go@golang.com",
  "password": "gog8g!pher"
}
```

response:
"OK"
_________________________________________________

#### Then we need to sign-in:

`POST /auth/sign-in`

```json
{
  "email": "go2g2o2@golang.com",
  "password": "aaaaaa"
}
```

response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjYyNzc4MjcsImlhdCI6MTY2NjE5MTQyNywic3ViIjoiMSJ9.WZeSW-tLPQZBXfhsTza_JLXKR9O2CGAgRh32uNzG1Tg"
}
```

_________________________________________________

#### Then we can create a post:

First, we need add authorization token in headers:
`Authorization: Bearer <token>`

Then we can create a post:
`POST /post`

```json
{
  "title": "web-develompent",
  "body": "THis is my first REST API in GO lang"
}
```

To get all posts:
`GET /post`

```json
[
  {
    "id": 1,
    "title": "web-develompent",
    "body": "THis is my first REST API in GO lang",
    "AuthorId": 1,
    "created_at": "2022-10-12T15:21:56.075473Z",
    "updated_at": "2022-10-12T15:21:56.075473Z"
  },
  {
    "id": 3,
    "title": "web-develompent 3",
    "body": "Slice in Go lang",
    "AuthorId": 1,
    "createdAt": "2022-10-12T15:25:31.069819Z",
    "updatedAt": "2022-10-12T15:25:31.069819Z"
  }
]
```
_________________________________________________

To get post by id:
`GET /post/<id>`

```json
  {
    "id": 1,
    "title": "web-develompent",
    "body": "THis is my first REST API in GO lang",
    "AuthorId": 1,
    "created_at": "2022-10-12T15:21:56.075473Z",
    "updated_at": "2022-10-12T15:21:56.075473Z"
  }
```
_________________________________________________

To update post by id:
`PUT /post`

```json
[
  {
    "id": 1,
    "title": "web-develompent UPDATED",
    "body": "THIS IS PYTHON",
  }
]
```
response:
```json
  {
    "message": "updated"
  }
```

_________________________________________________

### Swagger docs

to update need to run command: swag init -g cmd/main.go

Example of usage: http://localhost:8080/swagger/index.html
![img.png](swaggo_example.png)