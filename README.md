# portfolio

## Resources
### Mongodb driver
Documentation
```
https://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongodb client
```
go get go.mongodb.org/mongo-driver/mongo
```


## Docker
### Installing mongodb as a Docker container
```
docker run --name mongoPortfolio -p 27017:27017 -d mongo:latest
```

## TODO
- allow uploading files
- implement a template cache?
- make a nice header
- make a nice footer
- tidy up project structure
- store blog bodies as html