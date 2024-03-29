# alive mobile app

# purpose 
- this is the mobile application which is about tracking habit 
- the main purpose is personal development

# dependencies 

- firebase notificaiton
- flutter mobile UI
- golang backend 
- postgresql database 
- aws deployment(so soon)

## starting 

(for backend) 
- git clone project
- docker compose up 
- docker compose build 
- docker compose run --service-ports web bash 



## application folder structure 
## system design 


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
- go install github.com/cosmtrek/air@latest this is better building process on terminal. for more information. [read here](!https://github.com/cosmtrek/air) 
- go get gorm.io/gorm --> this is for database
-  go get gorm.io/driver/postgres --1 orm style makes manipulating simple in an object-oriented way 