new_migration:
	docker run --rm \
 		-v "$(shell pwd)/internal/script/migration:/db/migrations" \
 		--network host migrate/migrate \
 		create -ext sql -dir /db/migrations -seq $(MESSAGE_NAME)


up_migration:
	docker run --rm \
		-v "$(shell pwd)/internal/script/migration:/db/migrations" \
		--network host migrate/migrate \
		-path=/db/migrations \
		-database "mysql://root:mysql-db@tcp(localhost:3307)/social_media_app_db" up


down_migration:
	docker run -it --rm \
		-v "$(shell pwd)/internal/script/migration:/db/migrations" \
		--network host migrate/migrate \
		-path=/db/migrations \
		-database "mysql://root:mysql-db@tcp(localhost:3307)/social_media_app_db" down