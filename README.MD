






### how to create this project (step by step)
- added docker-compose.yml with the requirements for the project
- added dockerfile to create container 
- docker cpmpose up for running docker
- docker compose build for saving changes in the file 
- docker compose run --service-ports web bash --> this is to run file in docker server
- build every dependencies of project in docker
- go mod init github.com/damlaYasarr/aliveApp
- go get github.com/gofiber/fiber/v2
- go run cmd/main.go -b 0.0.0.0 -> starting server in docker