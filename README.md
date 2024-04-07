# social-network

## Инструкция по запуску:
1. установите docker и docker-compose.
2. в корне проекта создать файл `.env`
3. выполните команду: `make run`. после запуска данной команды поднимутся докер-контейнеры и накатятся миграции.
4. все необходимые команды лежат в [постман-коллекции](https://www.postman.com/aerospace-cosmonaut-29691174/workspace/highload-architect/collection/33337980-46a4c50d-5b28-4566-87dd-57e178216abd?action=share&creator=33337980)
5. (optional) - `make generate-data` генерирует и сохраняет в БД 1млн записей пользователей с реальными именами и фамилиями.