gen_authen_and_post:
	protoc --go_out=. --go-grpc_out=. internal/grpc/authen_and_post.proto

clean_authen_and_post:
	rm internal/grpc/pb/authen_and_post/*.go

gen_newsfeed:
	protoc --go_out=. --go-grpc_out=. internal/grpc/newsfeed.proto

clean_newsfeed:
	rm internal/grpc/pb/newsfeed/*.go

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