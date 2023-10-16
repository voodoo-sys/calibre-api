# calibre-api

### Overview

Dockerized http bridge to calibre utilities.

https://github.com/voodoo-sys/calibre-api

### Instructions


### Endpoints

* `POST /upload`  
    - Input
         - Description: multipart/form-data with "file" field containing file to upload
         - Example: `curl -F file=@test.pdf http://127.0.0.1:8080/upload`
    - Output
         - Description: json object with "filename" and "size" fields
         - Example: `{"filename":"1a8368fd-9aba-450a-b20f-d292d34fdfa0.pdf","filesize":9346436}`

* `POST /download`
    - Input
         - Description: json object with "filename" field
         - Example: `curl -X POST -H 'Content-Type: application/json' -d '{"filename":"1a8368fd-9aba-450a-b20f-d292d34fdfa0.pdf"}' http://127.0.0.1:8080/download`
    - Output
         - Description: binary data

* `GET /download/[:filename]`
    - Input
         - Description: path with "filename"
         - Example: `curl -X GET http://127.0.0.1:8080/download/1a8368fd-9aba-450a-b20f-d292d34fdfa0.pdf`
    - Output
         - Description: binary data

* `GET /files`
    - Input
         - Description: 
         - Example: `curl -X GET http://127.0.0.1:8080/files`
    - Output
         - Description: json object with array of json objects with "filename" and "size" fields
         - Example: `[{"filename":"1a8368fd-9aba-450a-b20f-d292d34fdfa0.pdf","filesize":9346436}]`

* `POST /ebook-convert`
    - Input
         - Description: json object with "filename", "format", "params" and "debug" field ("params" and "debug" are optional)
         - Example: `curl -X POST -H 'Content-Type: application/json' -d '{"filename":"1a8368fd-9aba-450a-b20f-d292d34fdfa0.pdf","format":"mobi","params":["--filter-css","font-family,color,margin-left,margin-right"]}' http://127.0.0.1:8080/ebook-convert`
    - Output
         - Description: json object with "filename" field
         - Example: `{"filename":"072c641a-738e-4d26-be14-039255f3299d.mobi"}`

* `GET /ebook-convert/version`
    - Input
         - Description: 
         - Example: `curl -X GET http://127.0.0.1:8080/ebook-convert/version`
    - Output
         - Description: json object with "version" field containing output of ebook-convert --version
         - Example: `[{"version":"ebook-convert (calibre 5.37.0)\nCreated by: Kovid Goyal \u003ckovid@kovidgoyal.net\u003e\n"}]`

### Environment

* PORT..
...Default: 8080..
...Description: Http server port...
...Example: `PORT=5000`..
