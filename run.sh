docker stop simpligo-pln
docker rm simpligo-pln
docker run -d --name simpligo-pln -p 80:8080 -p 443:443 -v $PWD:/shared -e SIMPLIGO_ENV=prod -e MAIN_SERVER_IP=127.0.0.1 -e PALAVRAS_IP=127.0.0.1 -e PALAVRAS_PORT=23080 --add-host elasticsearch:172.17.0.1 --restart always simpligo-pln:$1