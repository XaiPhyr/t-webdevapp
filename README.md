# Tutorial Web Development

uses [SQL Migrate](https://github.com/rubenv/sql-migrate)

## Folder struct

```
+---db
+---secrets
|   +---backend
|   \---db
\---src
    +---backend
    |   +---conf
    |   +---controllers
    |   +---middlewares
    |   +---models
    |   +---routers
    |   +---services
    |   +---sql
    |   |   \---migrations
    |   +---template
    |   |   +---css
    |   |   \---emails
    |   +---tests
    |   +---utils
    |   \---websocket
    \---frontend
        \---public
```

## Start development

### Secrets
 - configuration is stored in `secrets` folder,\
to initialize configration, rename template.env to `.env`

### Docker
- run command `docker compose up --watch` for windows
- edit html files inside `/frontend/public`