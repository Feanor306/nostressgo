server: 
	go run "src/main.go"

client:
	go run "src/client/main.go"

up: 
	docker-compose -f nostressgo.yml up -d

upb:
	docker-compose -f nostressgo.yml up -d --build --force-recreate

down:
	docker-compose -f nostressgo.yml down
	