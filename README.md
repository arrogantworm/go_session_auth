# Session Auth API

Session аутентификация на Go (cookie-based).

---

## Стек технологий

- **Go**
- **PostgreSQL**
- **Docker** + **docker-compose**
- **Swagger** (swaggo)

---

## Запуск проекта

```bash
docker compose -f docker-compose.yaml up --build
```

---

## Переменные окружения

```properties
PGUser=db_user
PGPassword=db_pass
PGHost=postgres
PGPort=5432
PGDBName=go_db
PGSSLMode=disable
```

---

## Конфигурация сессий

```yaml
sessions:
  TTL: 720h  # Время жизни сессии
  ```

---

## Документация API

```commandline
http://localhost:8000/swagger/index.html
```
