### Краткое описание работы
Реализован [микросервис нотификаций](https://github.com/syth0le/realtime-notification-service) в режиме реального времени
В работе использовался RabbitMQ в качестве очереди для общения между двумя микросервисами.
Также реализовано `internal api`, используя GRPC в качестве протокола. В апи вошел сервис Аутентификации, используемый для проверки валидности пользовательского токена, пришедшего в `realtime-notification-service`.

### Схема работы:
![notification.png](notification.png)

### Запуск и проверка
1. склонировать данный репозиторий и поднять его с помощью команды `make rebuild`
2. склонировать [микросервис нотификаций](https://github.com/syth0le/realtime-notification-service) и поднять его с помощью команды `make rebuild`
3. зарегистрироваться в ручке `/user/register` (для проверки будет достаточно двух пользователей).
4. залогиниться и получить свой токен для авторизации для обоих пользователей (в каждом последующем запросе его передавать в заголовке `Authorization`).
5. добавить в друзья пользователя-2, делая запрос от лица 1‑го пользователя.
6. создать несколько постов от обоих пользователей.
7. создать websocket коннект, используя ручку `/post/feed/posted` (передавая в заголовке `Authorization` токен пользователя) и проверить для каждого пользователя, что у них приходят обновления новостей.
8. все необходимые ручки есть в [постман коллекции](https://www.postman.com/aerospace-cosmonaut-29691174/workspace/highload-architect/collection/33337980-46a4c50d-5b28-4566-87dd-57e178216abd?action=share&creator=33337980). В работе можно использовать ручки в коллекции `5th hw`.

Особенности работы:
- учтено подключение нескольких устройств для одного пользователя. **Уведомление о посте прилетит на все подключенные устройства пользователя**.
- использован механизм Routing Key RabbitMQ в реализации.
- для масштабирования сервиса вебсокетов будет достаточно докинуть несколько инстансов (чтобы не перевалить за 65тыс возможных коннектов на один хост).
- К вебсокету можно без проблем подключиться и отключится. Серверная часть самостоятельно штатно завершит ненужные процессы, которые были созданы для работы нотификаций для пользователя.
- Процесс масштабирования RabbitMQ состоит следующим образом:
  - существует 2 этапа, которые могут помочь справиться с доступностью сервиса, с возможностями очереди обрабатывать поток сообщений.
  - Первый этап - горизонтальное масштабирование. Простое докидывание хостов
  - Второй этап - кластеризация и использование шардирования. Можно использовать механизм Hash Ring'а для определения в какой шард отправлять сообщения для определенного пользователя.
  - как вариант масштабирования - использование этого [инструмента](https://github.com/rabbitmq/rabbitmq-sharding)