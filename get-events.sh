#!/bin/bash

nc localhost 8080 << EOF
GET /events HTTP/1.1
Host: localhost
Connection: close

EOF
