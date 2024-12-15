## Run
```sh
cd ./cmd
go run main.go
```
___
## ENV

Default values
```ini
# database
POSTGRES_PASSWORD="postgres"
SMTP_SERVER_PASSWORD=""
```
___
## Endpoints
- **POST** `/sign-up`
- **POST** `/email/verify`
- **POST** `/email/refresh` 
 
*value see below*
___
## JSON DATA
*sign-up*
```json
{
  "first_name": "string",
  "last_name": "string",
  "username": "string",
  "email": "example@example.com",
  "password": "string"
}
```

*verify*
```json
{
  "email": "example@example.com",
  "code": "string[0-9]"
}
```
*refresh*
```json
{
  "email": "example@example.com",
}
```
___
## Technologies
**go packages**
- gin
- log/slog
- pgxpool
- viper

**utils**
- migrate