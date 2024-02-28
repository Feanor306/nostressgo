server: 
	go run "cmd/server/main.go"

client:
	go run "cmd/client/main.go"

up: 
	docker-compose -f nostressgo.yml up -d

upb:
	docker-compose -f nostressgo.yml up -d --build --force-recreate

down:
	docker-compose -f nostressgo.yml down
	