[![Build Status](https://travis-ci.org/byrnedo/apibase.svg?branch=master)](https://travis-ci.org/byrnedo/apibase)

# apibase

Base template/helper code for http apis.

Stack: 
 - Mongo
 - Nats mq
 - Protobuf

Skeleton code for some lightweight http services I want to write.

Current Functionality:

- Config
- Logger
- Various Middleware
- Docker test container helper
- Mongo wrapper - mgo used
- Postgres wrapper ( with migrations - gomigrate used )
- Nats MQ wrapper ( using protobuf ) nats-io/nats used
- Restclient

Template project to come.
