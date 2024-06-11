run:
	docker-compose up -d

rebuild:
	docker-compose up -d --build

# make make-migration ARGS="name"
make-migration:
	migrate create -ext sql -dir migrations -seq $(ARGS)

# migrate tool
migrate-master:
	migrate -path ./migrations -database "postgresql://social-network-local-admin:eephayl3eaph8Xo@localhost:6432/social-network-local?sslmode=disable" -verbose up

migrate-slave:
	migrate -path ./migrations -database "postgresql://social-network-local-admin:eephayl3eaph8Xo@localhost:7432/social-network-local?sslmode=disable" -verbose up

migrate:
	make migrate-master & make migrate-slave

# pgmigrate tool
migrate-master-pgmigrate:
	pgmigrate migrate -vvv -c "host=localhost port=6432 dbname=social-network-local user=social-network-local-admin password=eepha[l3eaph8Xo target_session_attrs=read-write sslmode=disable" -t latest

migrate-slave-pgmigrate:
	pgmigrate migrate -vvv -c "host=localhost port=7432 dbname=postgres user=postgres password=postgres target_session_attrs=read-write sslmode=disable" -t latest

migrate-pgmigrate:
	migrate-master-pgmigrate & migrate-slave-pgmigrate

generate-users:
	cd cmd/users-generator && go build && ./users-generator --config=local_config.yaml


#echo 'GET http://localhost:8080/user/search?first_name=na&second_name=na' | \                                                                                                                                          [±main ●]
#    vegeta attack -rate 5000 -duration 10s | vegeta encode | \
#    jaggr @count=rps \
#          hist\[100,200,300,400,500\]:code \
#          p25,p50,p95:latency \
#          sum:bytes_in \
#          sum:bytes_out | \
#    jplot rps+code.hist.100+code.hist.200+code.hist.300+code.hist.400+code.hist.500 \
#          latency.p95+latency.p50+latency.p25 \
#          bytes_in.sum+bytes_out.sum
