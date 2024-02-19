NOSTRESSGO_MAX_CONN := 5

server: 
	go run "src/main.go"

client:
	go run "src/client/main.go"