# URL Shortener


## Stack

- Golang
  - config: cleanenv
  - logger: slog
  - router: chi, "chi-render"
- DB
  - PostgresSQL

## Launch of the project

Do not store passwords in the project! There are secrets for this.


Add to environment values:

```yaml
CONFIG_PATH=config/local.yaml
DB_PORT=YOURS_PORT
DB_NAME=YOURS_DB_NAME
DB_USER=YOURS_DB_USER
DB_PASSWORD=YOURS_DB_PASSWORD
```