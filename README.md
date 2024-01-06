# sportsvoting

To run the project, you need to have MySQL installed on the system, after it is setup, and you know your MySQL username/password, run these commands in the terminal:
```
export DBUSER=<mysql-username>
export DBPASS=<mysql-password>
```

In MySQL create a database called 'nba':
```
CREATE DATABASE nba
```


After that, to have migrate library build the needed tables in the database, run:
```
migrate -path "./migrations/" -database "mysql://${DBUSER}:${DBPASS}@tcp(localhost:3306)/nba?multiStatements=true" up
```

(In case you want to delete these tables, just run the same command as above, just write 'down' instead of 'up' at the end of the command)

After that, you can run the project with a simple

```
go build
./sportsvoting
```

To run the frontend, in a separate terminal call:

```
npm run start
```

You can also run this app with docker compose with a simple call to:

```
docker compose up -d
```

Wait for it to start up, and then, find the IP address of the frontend service, enter it in the URL and the application should be up and running