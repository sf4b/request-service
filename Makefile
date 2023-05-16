build: go.sum
	go build -o ./main/main.go

go.sum:
	go mod tidy


run:
	go run ./main \
		-c ./main/config.template.yaml \
		-e ./main/env.sample
