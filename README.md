# VoluSnap - Cloud Volume Auto Snapshot Server

[![Build Status](https://travis-ci.org/didil/volusnap.svg?branch=master)](https://travis-ci.org/didil/volusnap)


*ALPHA: Not ready for production use*

Volusnap includes: 
- API server/daemon
- CLI client

## Usage

### Get Package
```
$ go get -u github.com/didil/volusnap
$ cd $GOPATH/src/github.com/didil/volusnap
$ make install
```

### PostgreSQL DB
Create a db and user by running the following commands in psql (replace name/password):
```
CREATE ROLE volusnap WITH LOGIN PASSWORD '123456';
CREATE DATABASE volusnap;
GRANT ALL PRIVILEGES ON DATABASE volusnap TO volusnap;
```

### Config files
```
$ cp sql-migrate.example.yml sql-migrate.yml
$ cp config.example.yml config.yml
```
Edit both file to adjust the configuration to your db credentials

### Build
```
$ make build
```

### Server
Start the Volusnap server:
```
$ ./volusnapd -p 8080
```

 ### Client
Signup:
```
$ ./volusnapctl signup -e "mike@example.com" -p "123456"
```

Login:
```
$ ./volusnapctl signup -e "mike@example.com" -p "123456"
INFO[0000] Login Successful Token:
xxJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoyLCJleHAiOjE1ODQ3OTExNjcsImlzcyI6ImFwcCJ9.C0JE7uh9SEL74ve53jKFkSh6fGZ5vIppXGOTymkRpdI 
```

Copy the JWT token to an env variable:
```
$ VOLUTOKEN=[token from previous command]
```

