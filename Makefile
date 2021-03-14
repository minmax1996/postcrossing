PWD = $(pwd)

docker-build:
	docker build -t minmax1996/postcrossing:latest .

docker-run:
	docker run -d -p 8080:8080 -v $(shell pwd)/postcards:/postcards/ minmax1996/postcrossing:latest

docker-run-withdiscord:
	docker run -d \
	-p 8080:8080 \
	-v $(shell pwd)/postcards:/postcards/ \
	-e "DISCORD_URL=yourwebhookurl" \
	minmax1996/postcrossing:latest