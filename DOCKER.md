
### Старт окружения
```bash
docker-compose up -d
```
Если docker новой версии, можно использовать `docker compose up -d`.

### Ребилд backend
```bash
docker-compose build backend
docker-compose up -d --no-deps backend
```

### Отключить все сервисы
```bash
docker-compose down
```

### Удалить окружение (вместе с базами данных)
```bash
docker-compose down --volumes --rmi all
```