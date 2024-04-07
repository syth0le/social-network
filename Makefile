run:
	docker-compose up -d

rebuild:
	docker-compose up -d --build

# make make-migration ARGS="name"
make-migration:
	migrate create -ext sql -dir migrations -seq $(ARGS)

# migrate tool
migrate-master:
	pgmigrate migrate -vvv -c "host=localhost port=6432 dbname=social-network-local user=social-network-local-admin password=eepha[l3eaph8Xo target_session_attrs=read-write sslmode=disable" -t latest

migrate-slave:
	pgmigrate migrate -vvv -c "host=localhost port=7432 dbname=postgres user=postgres password=postgres target_session_attrs=read-write sslmode=disable" -t latest

migrate:
	migrate-master & migrate-slave

# pgmigrate tool
migrate-master-pgmigrate:
	pgmigrate migrate -vvv -c "host=localhost port=6432 dbname=social-network-local user=social-network-local-admin password=eepha[l3eaph8Xo target_session_attrs=read-write sslmode=disable" -t latest

migrate-slave-pgmigrate:
	pgmigrate migrate -vvv -c "host=localhost port=7432 dbname=postgres user=postgres password=postgres target_session_attrs=read-write sslmode=disable" -t latest

migrate-pgmigrate:
	migrate-master-pgmigrate & migrate-slave-pgmigrate

generate-data:
	cd cmd/users-generator && go build && ./users-generator --config=local_config.yaml
