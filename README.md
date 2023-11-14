## Description

[Chat Realtime]

Tech Stack:

- Golanng (go fiber, gorilla/websocket)
- Postgres
- Supabase
- Kafka

## Install Migration

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## Run Service

```bash
$ docker-compose up -d
```

## Run Migration UP

```bash
$ make migration-up
```

## Run Migration Down

```bash
$ make migration-down
```

## Create Migration

```bash
$ make migration
#type your migration name example: create_create_table_users
```

## Test Coverage

```bash
$ make test-cov
```

## Register User (POST)

```bash
curl --location 'http://localhost:5001/users/registration' \
--header 'Content-Type: application/json' \
--data '{
    "name":"ahsan",
    "username":"ahsan",
    "password":"123"
}'
```

## Login User (POST)

```bash
curl --location 'http://localhost:5001/users/login' \
--header 'Content-Type: application/json' \
--data '{
    "username":"ahsan",
    "password":"123"
}'
```

## Send Message (PUT)

```bash
curl --location --request PUT 'http://localhost:5001/chat/send-message/{target_user}' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZXhwIjoxNzAwMDIwNDg3fQ.IA5Oz4Un-_rpr9SNjif4bX6flfhP43dViAg_kdjBRYk' \
--header 'Content-Type: application/json' \
--data '{
    "message":"haloo.."
}'
```

## Fetch Message (GET)

```bash
curl --location 'http://localhost:5001/chat/{chat_target_user}' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZXhwIjoxNzAwMDIwNDg3fQ.IA5Oz4Un-_rpr9SNjif4bX6flfhP43dViAg_kdjBRYk'
```

## Search Message (GET)

```bash
curl --location 'http://localhost:5001/chat/search?message=halo' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZXhwIjoxNzAwMDI1OTAwfQ.q4MeOCp8aZfGGYniOZc79ElRl5pcKkNOu7FPx0fUR3k'
```

## Fiber Monitoring

Available at `http://localhost:5001`

## Flowchart

![alt text](https://github.com/dedihartono801/chat-realtime/blob/master/flowchart-send-message.png)

## ERD

![alt text](https://github.com/dedihartono801/chat-realtime/blob/master/ERD.png)
