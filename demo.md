dapr run --app-id myapp --dapr-http-port 3500

second terminal:

curl -X POST -H "Content-Type: application/json" -d '[{ "key": "name", "value": "Bruce Wayne"}]' http://localhost:3500/v1.0/state/statestore


curl http://localhost:3500/v1.0/state/statestore/name


docker exec -it dapr_redis redis-cli

hgetall "myapp||name"



third ternminal:

curl -X POST -H "Content-Type: application/json" -d '[{ "key": "name", "value": "Robert Wayne"}]' http://localhost:3500/v1.0/state/statestore


Terminal 2:

hgetall "myapp||name"

set name "Michael Jordan"

exit


Terminal 3: 

curl http://localhost:3500/v1.0/state/statestore/name


================
MySQL

Terminal 2: 

run a MySQL Database

docker run --name dapr-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:latest

To connect to the MySQL server from the MySQL client application from inside the container running MySQL:

docker exec -it dapr-mysql mysql -uroot -p

and type the password: `my-secret-pw`

show databases;


Ternminal 3 

Create a file called `mysql-statestore.yaml` and add the following content to the file:

apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: durable-statestore
spec:
  type: state.mysql
  version: v1
  metadata:
  - name: connectionString
    value: "root:my-secret-pw@tcp(localhost:3306)/?allowNativePasswords=true"


Then run a Dapr sidecar from the directory that contains file `mysql-statestore.yaml`, as follows:

dapr run --app-id myotherapp --dapr-http-port 3510 --components-path .


In terminal 2 (MySQL client):
show databases;
use dapr_state_store;
show tables;
select * from state;

Terminal 4:

curl -X POST -H "Content-Type: application/json" -d '[{ "key": "name", "value": "Your Own Name"}]' http://localhost:3510/v1.0/state/durable-statestore

In terminal 2 (MySQL client):
select * from state;


Terminal 4:

curl http://localhost:3510/v1.0/state/durable-statestore/name

Terminal 1, 3:
stop Dapr Sidecar

Terminal 4:

curl http://localhost:3510/v1.0/state/durable-statestore/name


Terminal 3:
dapr run --app-id myotherapp --dapr-http-port 3510 --components-path .


Terminal 4:

curl http://localhost:3510/v1.0/state/durable-statestore/name
