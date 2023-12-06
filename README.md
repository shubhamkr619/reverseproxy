# go-revproxy

## Intro
This is the code base of reverse proxy, this allows us to configure upstream using mysql database.
This is written in golang to support only CDN with gcs bucket as upstream while allowing to change bucket based on the hostname. 

## Config 
Mysql-connection property and Application ports will be required to start the application. Sample config file to placed as `config.yml` in config folder to make use of this.

```yaml
server:

  Port: 8080

database:
  dbName: ""
  dbUser: ""
  dbPass: ""
  dbHost: ""
  dbPort: 3306


```

## Example

```shell
curl -X GET localhost:8080/Daytona.jpeg -H 'Host: customer2.com'  -I

HTTP/1.1 200 OK
Accept-Ranges: bytes
Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000
Cache-Control: public, max-age=3600
Content-Length: 121884
Content-Type: image/jpeg
Date: Wed, 06 Dec 2023 12:25:49 GMT
Defaultexpire: 300
Etag: "00c34c812f401ba94086fd48d8b55b86"
Expires: Wed, 06 Dec 2023 13:25:49 GMT
Last-Modified: Thu, 05 Oct 2023 13:47:11 GMT
Server: UploadServer
X-Goog-Generation: 1696513631229626
X-Goog-Hash: crc32c=QQFILw==
X-Goog-Hash: md5=AMNMgS9AG6lAhv1I2LVbhg==
X-Goog-Metageneration: 1
X-Goog-Storage-Class: STANDARD
X-Goog-Stored-Content-Encoding: identity
X-Goog-Stored-Content-Length: 121884
X-Guploader-Uploadid: ABPtcPorhrWMVx1vc_khcy91EhJPRfFcGlUO3BXewRRt4o4IKDoZbTT8hpo-QiJYvwpYIbnodg82RXRJNQ
```

### Failure 
```shell
curl -X GET localhost:8080/Daytona.jpeg -H 'Host: customer3.com'  -I

HTTP/1.1 403 Forbidden
Cache-Control: private, max-age=0
Content-Length: 298
Content-Type: application/xml; charset=UTF-8
Date: Wed, 06 Dec 2023 12:25:39 GMT
Defaultexpire: 300
Expires: Wed, 06 Dec 2023 12:25:39 GMT
Server: UploadServer
X-Guploader-Uploadid: ABPtcPpot9MljypFar_a7mkLNE_v0cFD6pi1ewmh7V0zMZJ2snZllnYHGRmhLc2Bx8FEkXlYewDgYfJRiA
```