## Реализация 
- в [docker-compose.yaml](https://github.com/syth0le/social-network/blob/main/docker-compose.yaml) добавлен контейнер со слейвами и миграторы для них

## Алгоритм пошаговый:
1. поднимаем приложение `make rebuild`
2. заходим в мастер `su - postgres`
3. cоздаем пользователя: `createuser --replication -P repluser' и задаем пароль `pass`
4. добавляем в файл `/data/master/postgresql.conf` следующее:
    ```
    wal_level = replica
    max_wal_senders = 2
    max_replication_slots = 2
    hot_standby = on
    hot_standby_feedback = on
    ```
5. для файла `/data/postgresql_01/pg_hba.conf` следующее:
   `host    replication     all             Subnet-IN-DOCKER           md5`
6. и рестартим все make run
7. в слейве выполняем команду: `su - postgres -c "pg_basebackup --host=master --username=repluser --pgdata=/var/lib/postgresql/data --wal-method=stream --write-recovery-conf"`
8. проверка работы мастера: `docker exec -it master su - postgres -c "psql -c 'select * from pg_stat_replication;'"`
9. проверка работы слейва: `docker exec -it slave su - postgres -c "psql -c 'select * from pg_stat_wal_receiver;'""`
10. можно добавить в postgresql.conf: `primary_conninfo = 'host=msater port=5432 user=repluser password=pass application_name=slave'`
11. проводим тестирование с синхронной репликацией
12. теперь перейдем к асинхронной реплике
13. настраиваем слейв 2 и включаем синхронную репликацию ```sql
    synchronous_commit = on
    synchronous_standby_names = 'FIRST 1 (pgslave, pgasyncslave)'

   select pg_reload_conf();
   ```
14. проводим нагрузочное тестирование снова
    `bombardier -c 100 -d 10s "localhost:8080/user/search-tarantool?first_name=Абр&second_name=Юр"`
15. убиваем мастер
16. промоутим слейв:
    ```
    docker exec -it slave su - postgres -c psql

    select * from pg_promote();
    
    synchronous_commit = on
    synchronous_standby_names = 'ANY 1 (master, slave-2)'
    ```
    
17. Подключим вторую реплику к новому мастеру 
    `primary_conninfo = 'host=slave port=5432 user=replicator password=pass application_name=slave-2'`

18. Восстановим мастер в качестве реплики
    ```
    touch master/standby.signal
    primary_conninfo = 'host=slave port=5432 user=replicator password=pass application_name=master'
    ```
19. как итог транзакции не потерялись и все доехали до асинхронной реплики