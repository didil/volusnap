# VoluSnap - Cloud Volume auto-snapshot Go Server

[![Build Status](https://travis-ci.org/didil/volusnap.svg?branch=master)](https://travis-ci.org/didil/volusnap)

Volusnap allows triggering automated recurring snapshots of cloud provider volumes. Digital Ocean and Scaleway APIs are currently supported.  

*ALPHA: Not ready for production use*

## Contributing 
Please open PRs to add providers !  
The architecture is modular, you can look at a provider example [here](pkg/api/digitalocean.go). Don't forget to add tests :)


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
INFO[0000] Starting snapRulesChecker ...                
INFO[0000] Starting server on port 8080 ...    
```

## API
The server uses a REST API. All POST requests expect "application/JSON" content-type. Responses are JSON

### Signup user
POST /api/v1/auth/signup/
*Request Body*
```json
{
    "email": "myemail@example.com",
    "password": "mypassword"
}
```
*Response Body*
```json
{
  "id": 1
}
```

### Login user  
POST /api/v1/auth/login/  
*Request Body*
```json
{
    "email": "myemail@example.com",
    "password": "mypassword"
}
```
*Response Body*
```json
{
  "token": "xxxxxxx"
}
```

### Create account
POST /api/v1/auth/account/  
*Headers*  
Authorization: "Bearer [Token From Login]"  
*Request Body*
```json
{
	"provider": "scaleway",
	"name": "scaleway 1", 
	"token": "xyzxyzxyz" // API Token from the cloud provider
}
```
*Response Body*
```json
{
  "id": 2
}
``` 

### List accounts 
GET /api/v1/auth/account/ 
Headers  
Authorization: "Bearer [Token From Login]"  
*Response Body*
```json
{
  "accounts": [
    {
      "id": 1,
      "created_at": "2019-02-25T18:38:55.436512+07:00",
      "updated_at": "2019-02-25T18:38:55.436512+07:00",
      "name": "do 1",
      "provider": "digital_ocean",
      "token": "xxx",
      "user_id": 1
    },
    {
      "id": 2,
      "created_at": "2019-02-26T14:35:05.261556+07:00",
      "updated_at": "2019-02-26T14:35:05.261556+07:00",
      "name": "scaleway 1",
      "provider": "scaleway",
      "token": "xxx",
      "user_id": 1
    }
  ]
}
``` 

### List Volumes
GET /api/v1/account/{AccountID}/volume/  
*Headers*  
Authorization: "Bearer [Token From Login]"  
*Response Body*
```json
{
  "volumes": [
    {
      "ID": "MY-VOLUME-ID",
      "Name": "MY-VOLUME-NAME",
      "Size": 50,
      "Region": "MY-VOLUME-REGION"
    },  
  ]
}
```

### Create SnapRule
POST /api/v1/account/{AccountID}/snaprule/  
*Headers*  
Authorization: "Bearer [Token From Login]"  
*Request Body*
```json
{
	"frequency": 168, // Frequency in hours
	"volume_id": "MY-VOLUME-ID",
	"volume_name": "MY-VOLUME-NAME",
	"volume_region": "MY-VOLUME-REGION"
}
```
*Response Body*
```json
{
  "id": 1
}
```

### List SnapRules
GET /api/v1/account/{AccountID}/snaprule/  
*Headers*  
Authorization: "Bearer [Token From Login]"  
*Response Body*
```json
{
  "snaprules": [
    {
      "id": 1,
      "created_at": "2019-02-26T14:14:44.194276+07:00",
      "updated_at": "2019-02-26T14:14:44.194276+07:00",
      "frequency": 168,
      "volume_id": "MY-VOLUME-ID",
      "volume_name": "MY-VOLUME-NAME",
      "volume_region": "MY-VOLUME-REGION",
      "account_id": 1
    }
  ]
}
```


### List Snapshots
GET /api/v1/account/{AccountID}/snapshot/  
*Headers*  
Authorization: "Bearer [Token From Login]"  
*Response Body*
```json
{
  "snapshots": [
    {
      "id": 1,
      "created_at": "2019-02-26T19:22:06.094935+07:00",
      "updated_at": "2019-02-26T19:22:06.094935+07:00",
      "provider_snapshot_id": "1234567",
      "snap_rule_id": 1
    }
  ]
}
```