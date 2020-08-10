heroku_registry = registry.heroku.com/go-shopify-api/web
binary_name = go_shop_api

all: debug

debug: build
	docker-compose up

run: build
	docker-compose up -d

build: 
	docker-compose build

heroku_deploy: clean
	sudo docker image build -t $(heroku_registry) .
	docker push $(heroku_registry)
	heroku container:release web

clean:
	rm -f ./server/$(binary_name)
	docker-compose down
	docker rmi -f $(heroku_registry)
	docker rmi -f $(binary_name)
	docker image prune -f --filter label=stage=builder