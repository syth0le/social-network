## Цель Реализовать функционал:
Вынести систему диалогов в отдельный сервис.

### System Design
![dialog.png](dialog.png)

В ходе реализации были проведены следующие работы:
- реализован микросервис [диалогов](https://github.com/syth0le/dialog-service)
  - вынесено хранение диалогов и сообщений в отдельную обособленную от монолита БД
  - реализовано [Internal GRPC API](https://github.com/syth0le/dialog-service/blob/main/proto/internalapi/dialog_service.proto) в которое ходит монолит, когда в его старую апишку пришел запрос
- логика старой апи
  - в сервисе монолите переведена на походы в новый сервис [link](https://github.com/syth0le/social-network/blob/main/internal/handler/publicapi/dialog.go)
  - старые клиенты могут все еще ходить в старое апи и не ломаться

TODO:
- сквозная регистрация запросов (сделать можно через x-request-id хэдер прокидываемый в каждом запросе)

### Запуск и проверка
1. склонировать данный репозиторий и поднять его с помощью команды `make rebuild`
2. склонировать [диалогов](https://github.com/syth0le/dialog-service) и поднять его с помощью команды `make rebuild`
3. зарегистрироваться в ручке `/user/register` (для проверки будет достаточно двух пользователей).
4. залогиниться и получить свой токен для авторизации для обоих пользователей (в каждом последующем запросе его передавать в заголовке `Authorization`).
5. создать диалог `[POST] localhost:8070/dialog` (в ручке сервиса диалогов) или в `[POST] localhost:8080/dialog`.
6. отправить сообщение в диалог `[POST] localhost:8070/dialog/send` (в ручке сервиса диалогов) или в `[POST] localhost:8080/dialog/send`.
7. получить все сообщения в диалоге `[GET] localhost:8070/dialog/{dialogId}/list` (в ручке сервиса диалогов) или в `[GET] localhost:8080/dialog/{dialogId}/list`.
8. все необходимые ручки есть в [постман коллекции](https://www.postman.com/aerospace-cosmonaut-29691174/workspace/highload-architect/collection/33337980-46a4c50d-5b28-4566-87dd-57e178216abd?action=share&creator=33337980). В работе можно использовать ручки в коллекции `8th hw http collection`.
