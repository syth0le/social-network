# Индексы и оптимизация запросов:

## При выполнении работы было реализовано: 

- генератор пользователей на Го [тут](../../cmd/users-generator), который делает insert в базу батчом по 1000шт
- сгенерировано 5828459 анкет пользователей
- добавлена ручка /search из спецификации

## Графики:

###  10 rps

##### latency до индекса
![img.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img.png)
##### throughput до индекса
![img_1.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_1.png)
##### http codes до индекса
![img_2.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_2.png)
##### latency после индекса
![img_9.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_9.png)
##### throughput после индекса
![img_10.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_10.png)
##### http codes после индекса
![img_11.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_11.png)

###  100 rps

##### latency до индекса
![img_3.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_3.png)
##### throughput до индекса
![img_4.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_4.png)
##### http codes до индекса
![img_5.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_5.png)
##### latency после индекса
![img_12.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_12.png)
##### throughput после индекса
![img_13.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_13.png)
##### http codes после индекса
![img_14.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_14.png)

###  1000 rps

##### latency до индекса
![img_6.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_6.png)
##### throughput до индекса
![img_7.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_7.png)
##### http codes до индекса
![img_8.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_8.png)
##### latency после индекса
![img_15.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_15.png)
##### throughput после индекса
![img_16.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_16.png)
##### http codes после индекса
![img_17.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_17.png)

###  10000 rps

##### без индекса - приложение может выдержать больше, но умирает база, не выдерживает такое количество коннектов.
##### А вот с индексом все работает замечательно. Добавление 1го доп пода приложения решит проблему небольшого фона пятисоток.

##### latency после индекса
![img_18.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_18.png)
##### throughput после индекса
![img_19.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_19.png)
##### http codes после индекса
![img_20.png](https://github.com/syth0le/social-network/blob/main/docs/2hw/img_20.png)


## По пунктам:
- запрос добавления индекса
    ```sql
    CREATE EXTENSION IF NOT EXISTS pg_trgm;
    CREATE EXTENSION IF NOT EXISTS btree_gin;
    CREATE INDEX user_first_name_idx ON user_table USING gin (first_name gin_trgm_ops);
    CREATE INDEX user_second_name_idx ON user_table USING gin (second_name gin_trgm_ops);
    ```

- explain запросов до индекса
```
Gather  (cost=1000.00..173972.07 rows=2 width=146)
  Workers Planned: 2
  ->  Parallel Seq Scan on user_table  (cost=0.00..172971.87 rows=1 width=146)
        Filter: ((deleted_at IS NULL) AND (first_name ~~* 'na%'::text) AND (second_name ~~* 'na%'::text))
JIT:
  Functions: 4
"  Options: Inlining false, Optimization false, Expressions true, Deforming true"
```

- explain запросов после индекса
```
Bitmap Heap Scan on user_table  (cost=410.54..418.55 rows=2 width=146)
  Recheck Cond: ((first_name ~~* 'na%'::text) AND (second_name ~~* 'na%'::text))
  Filter: (deleted_at IS NULL)
  ->  BitmapAnd  (cost=410.54..410.54 rows=2 width=0)
        ->  Bitmap Index Scan on user_first_name_idx  (cost=0.00..68.62 rows=352 width=0)
              Index Cond: (first_name ~~* 'na%'::text)
        ->  Bitmap Index Scan on user_second_name_idx  (cost=0.00..341.67 rows=36040 width=0)
              Index Cond: (second_name ~~* 'na%'::text)
```

- использовал модуль pg_trgm так как он позволяет находить похожие строки, даже если они не совпадают буквенно. Это полезно, например, при исправлении опечаток или при поиске синонимов.
- Трехграммы (pg_trgm) учитывают порядок символов и поддерживают операторы для сравнения, что делает их мощным инструментом для полнотекстового поиска.
- Выбирал индекс из двух: GIN или GIST. Остановился на GIN так как он выполняет поиск по индексу примерно в 3 раза быстрее чем GIST.
- Даже несмотря на то, что пересчет индекса GIN наиболее ресурсозатратная операция, был выбран именно он, так как добавление и регистрация новых пользователей происходит на порядок реже чем чтение информации о них.
- Добавление индекса на отдельные колонки позволяет гибко искать пользователя по имени и фамилии.