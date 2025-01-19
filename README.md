# Historical_Publications_Back

## КАК ЗАПУСТИТЬ

1) #### Запускаем МИНИО через докер (заранее подготовить изображения для миграции в папке migrations/events)
```
docker-compose up -d
```
2) #### Заходим на http://localhost:9001/buckets/events/admin/summary - там ставим Access Policy на Public
3) #### Запускаем миграцию БД
```
go run ./cmd/migrate/main.go  
```
4) #### Запускаем сервер 
```
go run ./cmd/main.go  
```

