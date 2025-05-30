# MySQLコンテナの起動
start:
	docker-compose up -d

# MySQLコンテナの停止
stop:
	docker-compose down

# MySQLコンテナの削除とボリュームの削除
clean:
	docker-compose down -v
	rm -rf ./db_data/*

# MySQLログを確認
logs:
	docker-compose logs -f db

# MySQLに接続
mysql:
	docker exec -it my_db mysql -u root -p

# appを実行
run:
	(cd app; go run main.go;)

# skeemaを実行
# my_db: DIR=my_db DB=my_db
# sakila: DIR=sakila DB=sakila
# migrate: DIR=migrate DB=migrate
# ENV例: ENV="--environment=local"
skeema-init:
	skeema init -h 127.0.0.1 -P 3306 -u root -proot -d schemas/table/${DIR} --schema ${DB};

skeema-diff:
	(cd schemas/table/${DIR}/; \
	date >> .skeema.log; \
	eval skeema diff --allow-unsafe ${ENV} >> .skeema.log;)

skeema-push:
	(cd schemas/table/${DIR}/; \
	date >> .skeema.log; \
	eval skeema push --allow-unsafe ${ENV} >> .skeema.log;)

skeema-pull:
	(cd schemas/table/${DIR}/; \
	date >> .skeema.log; \
	eval skeema pull -proot ${ENV} >> .skeema.log;)

# golang-migrateを実行
# migrate: DIR=migrate DB=migrate
migrate-create:
	migrate create -ext sql -dir ./schemas/migration/${DIR} -seq -digits 6 ${NAME}

migrate-up:
	migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/${DB}" -path ./schemas/migration/${DIR} up

migrate-down:
	migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/${DB}" -path ./schemas/migration/${DIR} down

migrate-drop:
	migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/${DB}" -path ./schemas/migration/${DIR} drop

# sakilaデータのインポート
import-sakila:
	docker cp ./sakila/sakila-schema.sql my_db:/sakila-schema.sql
	docker cp ./sakila/sakila-data.sql my_db:/sakila-data.sql
	docker exec -i my_db mysql -u root -proot my_db < ./sakila/sakila-schema.sql
	docker exec -i my_db mysql -u root -proot my_db < ./sakila/sakila-data.sql
