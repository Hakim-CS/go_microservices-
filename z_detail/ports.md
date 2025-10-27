# Test broker
curl http://localhost:8080/ping

# Test authentication
curl http://localhost:8081/ping

# Test logger
curl http://localhost:8082/ping

# Test mail
curl http://localhost:8083/ping

# MongoDB
mongodb://localhost:27017

# postgres
mongodb://localhost:5432


# mongoDB connection query:
mongodb://admin:password@localhost:27017/logs?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false