# Hotel Reservation Backend

## Outline

- users -> book a room from an hotel
- admins -> check reservations
- Hotels -> CRUD API -> JSON
- Rooms -> CRUD API -> JSON
- Authentication and Authorization -> JWT tokens  
- sctipts for database management -> seeding -> migrations


## Install mongodb as a docker container



- [Install MongoDB Community with Docker](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-community-with-docker/)

```bash
docker run --name mongo   -p 27017:27017  -d mongodb/mongodb-community-server:latest
```

**.env**
```
DB_URI=mongodb://localhost:27017
```



- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)

 