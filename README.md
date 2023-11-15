# howlongtobeatmybacklog

## Description 

To build with the list of game you own a backlog of "To want to play" and estimate how long it would take.
Go as a backend to brush up my skills in it, and rust with Yew in the frontend to learn.

## How to run

### Backend

`cd go`
`docker compose up`

### Frontend

`cd rust`  
`docker build  -t myrust . `  
`docker run --rm -p 127.0.0.1:8000:8000 myrust`  

## MVP

- [ ] Have a go backend that allows to store owned games from a steam id
- [ ] This go backend then serves for a FE this list
- [ ] The frontend needs to be able to query that list and display it
- [ ] Being able to scrape or with APIS, the howlongtobeat.com to estimate length of the game
- [ ] In the frontend pick some games and generate a backlog with a full estimate on how long to complete


## Strech goals

- Search, filtering, sorting of the list of games
- set all in GCP

### Backend
- https://github.com/mongodb/mongo-go-driver#network-compression
- set pooling connection for mongo

### Frontend
- with the time estimation 

## References
- https://www.linkedin.com/pulse/exciting-golang-mongodb-tips-aditira-jamhuri/
- https://www.mongodb.com/docs/drivers/go/current/usage-examples/
- https://earthly.dev/blog/use-mongo-with-go/
- https://earthly.dev/blog/mongodb-docker/
- https://www.mongodb.com/developer/languages/go/get-hyped-using-docker-go-mongodb/
- https://go-recipes.dev/parsing-json-with-go-7268937a5f7b
- https://yew.rs/docs/tutorial#fetching-data-using-external-rest-api
