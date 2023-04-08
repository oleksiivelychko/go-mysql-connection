docker-network:
	docker network inspect go-network >/dev/null 2>&1 || docker network create --driver bridge go-network

docker-volume:
	docker volume inspect mysql-data >/dev/null 2>&1 || docker volume create mysql-data

mysql-client-run: docker-network
	docker run -it --network go-network mysql:8.0 mysql -hmysql-server -uroot -p

mysql-server-run: docker-network docker-volume
	docker run --name mysql-server \
		--network go-network \
		-v mysql-data:/var/lib/mysql \
		-p 3306:3306 \
		-e MYSQL_ROOT_PASSWORD=secret \
		mysql:8.0

mysql-stop-force:
	docker stop mysql-server
	docker rm mysql-server
