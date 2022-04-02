# Handson for Conclusion Code Cafe on Dapr.io - April 2022

- [Introduction to Dapr.io and](#introduction-to-daprio-and)
  - [1. Setting up Dapr.io](#1-setting-up-daprio)
    - [a. Installing the Dapr.io CLI](#a-installing-the-daprio-cli)
    - [b. Initialize Dapr in your local environment](#b-initialize-dapr-in-your-local-environment)
  - [2. Playing with Dapr State Store capability](#2-playing-with-dapr-state-store-capability)
  - [3. Using a MySQL Database as a Dapr State Store](#3-using-a-mysql-database-as-a-dapr-state-store)
  - [4. Telemetry](#4-telemetry)
  - [Closure](#closure)
  - [Resources](#resources)

In this lab, you will install Dapr.io in your local environment and explore some of its capabilities as a *personal assistant* for applications and microservices. Dapr.io offers a variety of services to applications, including state management, asynchronous communication (pub & sub), handling secrets and configuration data, protection, routing and load balancing of communication and interacting with external services and technology platforms. Applications rely on Dapr.io to handle *dirty details* such as technology specific APIs and configuration - and only talk to Dapr.io in civilized terms - a standardized, functional API.  

## 1. Setting up Dapr.io

Dapr.io runs very well on Kubernetes. However, for simplicity sake we will go for the more straightforward approach. You will need an environment - MacOS, Windows, Linux - that has Docker running.  

Note: there is the possibility to run Dapr.io without Docker at all - see [Dapr.io self hosted without Docker](https://docs.dapr.io/operations/hosting/self-hosted/self-hosted-no-docker/) for details; however, that requires a lot of additional work and will limit your progress through today's hands on labs.

### a. Installing the Dapr.io CLI

Follow the instructions for installing the Dapr.io CLI that are provided on [this page](https://docs.dapr.io/getting-started/install-dapr-cli/).

Using
```
dapr
```
you should now have feedback that indicates a properly installed Dapr environment.

### b. Initialize Dapr in your local environment

Follow the instructions for initializing the Dapr environment that are provided on [this page](https://docs.dapr.io/getting-started/install-dapr-selfhost/). The outcome of this step should be three running containers: Dapr, Zipkin and Redis. The Redis container provides the default implementation for State Store and Pub Sub broker used by Dapr based on the Redis in memory cache. Zipkin collects telemetry and provides insight in tracing of requests to and from microservices through Dapr.

At this point, you have a virtual pool of personal assistants to draw support from for each of the applications and microservices you will run. These potential assistants when called into action can make use of the standard Dapr facilities for state management and pub/sub and for collecting telemetry and routing requests. We can easily add - and will so later in this lab - additional capabilities that these PAs can make use of.  

![](images/standalone-dapr-setup.png)
## 2. Playing with Dapr State Store capability

Follow the instructions for trying out the default Dapr state store implementation (based on Redis cache) as described in this document [Making the Dapr PA hold and return state ](https://docs.dapr.io/getting-started/get-started-api/). 

In these steps, you run a Dapr Sidecar (the technical term for the personal assistant) for an application called *myapp*. That application does not actually exist. What we do is hire a personal assistant for a manager who has not arrived yet. The PA can start doing their work - but not yet pass messages on to the person they are the PA for. In this case, we can ask the sidecar for example to save state and publish messages - but any attempt to invoke the *myapp* application will fail because it is not there yet. 

When you check the logs after executing the first command, you will see messages indicating that the sidecar (again, the technical term for the personal assstant) is running. Because we did not explicitly configure components for state store and pub sub, the default implementations in Dapr are used - based on Redis and leveraging the Docker Container running the Redis image.

After completing the steps in that document, execute the following command:
```
curl -X POST -H "Content-Type: application/json" -d '[{ "key": "name", "value": "Your Own Name"}]' http://localhost:3500/v1.0/state/statestore
```
This sets a new value for the state stored under the key *name*, updating the value stored previously under that key. Note that in the URL we refer to localhost and the port on which the sidecar was started through the *dapr run* command. The path segment *v1.0* refers to the version of the Dapr APIs. The segment *state* indicates that we are invoking the State Store API. Finally, the URL path segment *statestore* indicates the name of the state store that we want to save state to. A Daprized application - or rather its Dapr sidecar - can work with multple state stores that each have their own name and can each have their own implementation. 

## 3. Using a MySQL Database as a Dapr State Store

As was discussed before, Dapr.io supports many different technologies. It provides for example over a dozen implementations of the state store component - each leveraging a different database or cloud storage service. In the previous section we used the out of the Dapr.io box provided Redis Cache as state store. In this step we are going to use a MySQL Database to serve the same purpose.

To run a MySQL Database

```
docker run --name dapr-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:latest
```

To connect to the MySQL server from the MySQL client application from inside the container running MySQL:
```
docker exec -it dapr-mysql mysql -uroot -p
```
and type the password: `my-secret-pw`

You can list the databases:
```
show databases;
```
and create a database and then create tables if you want to.

However, let's not do that right now. Let's configure Dapr to use this MySQL instance to create a database and table to store state in.

Details can be found in the [documentation on the Dapr MySQL State Store building block](https://docs.dapr.io/reference/components-reference/supported-state-stores/setup-mysql/).

Start a new terminal session. Create a file called `mysql-statestore.yaml` and add the following content to the file:

```
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
```

Then run a Dapr sidecar from the directory that contains file `mysql-statestore.yaml`, as follows:
```
dapr run --app-id myotherapp --dapr-http-port 3510 --components-path .
```
Note: if the current directory contains other yaml-files you may see unexpected and unintended effects as Dapr tries to interpret them as well.

This instruction starts a Dapr sidecar (a personal assistant) and instructs the sidecar about a state store called *durable-statestore*. This statestore is backed by a MySQL Database for which the connection details are provided. Now when anyone asks this sidecar to save state and specifies the *durable-statestore* as the state store to use for that request, the Dapr sidecar will know where to go and because of the many built in building blocks in Dapr it also knows what to do in order to talk state affairs with MySQL.

This is what your environment looks like at this point:
![](images/mysql-statestore.png)

You will find lines like these ones in the logging produced by Dapr when starting up:
```
INFO[0000] Creating MySql schema 'dapr_state_store'      app_id=myotherapp instance=DESKTOP-NIQR4P9 scope=dapr.contrib type=log ver=edge
INFO[0000] Creating MySql state table 'state'            app_id=myotherapp instance=DESKTOP-NIQR4P9 scope=dapr.contrib type=log ver=edge
INFO[0000] component loaded. name: durable-statestore, type: state.mysql/v1  app_id=myotherapp instance=DESKTOP-NIQR4P9 scope=dapr.runtime type=log ver=edge
```

This confirms that Dapr initialized communications with the MySQL instance, it also created the default schema and default table in it for storing state.

Let us now create some state, in exactly the same way as we created state before - when it was saved in Redis Cache.

```
curl -X POST -H "Content-Type: application/json" -d '[{ "key": "name", "value": "Your Own Name"}]' http://localhost:3510/v1.0/state/durable-statestore
```
Note that the portname at which we access the Dapr sidecar is 3510 and the name of the statestore requested is passed in the URL path as well. Let's check if the state was saved. First by retrieving it from the sidecar:
```
curl http://localhost:3510/v1.0/state/durable-statestore/name
```

And next by checking directly in the MySQL Database.
To connect to the MySQL server as you did before, run this next statement that opens the MySQL client in the container running MySQL:
```
docker exec -it dapr-mysql mysql -uroot -p
```
and type the password: `my-secret-pw`

You can list the databases:
```
show databases;
```
and you will notice a database called *dapr_state_store* has been created. 

Use these next statements to switchh to the *dapr_state_store* database, to list all tables and to select all records from the one table *STATE* that Dapr created when it initialized the *durable-statestore* component.

```
use dapr_state_store;
show tables;
select * from state;
```
The last statement returns a result similar to this one:
```
+------------------+-----------------+----------+---------------------+---------------------+--------------------------------------+
| id               | value           | isbinary | insertDate          | updateDate          | eTag                                 |
+------------------+-----------------+----------+---------------------+---------------------+--------------------------------------+
| myotherapp||name | "Your Own Name" |        0 | 2022-03-01 18:00:01 | 2022-03-01 18:00:01 | 0a4d1bb3-e208-4c03-8296-eb9f1a544ff3 |
+------------------+-----------------+----------+---------------------+---------------------+--------------------------------------+
```
Note how the *key* in column *id* is composed of two parts: the name of the application for which the sidecar was managing state concatenated to the actual key you specified. 

The state held in table *STATE* should typically be managed only through Dapr. However, if you were to change the state directly through SQL:
```
update state set value = '"Somebody Else Altogether"';
```
You can exit the mysql client by typing `exit`. This returns you to the command prompt.

You will find that when you ask Dapr for the state held under key *name* it will return the updated value, once again proving that Dapr interacts with MySQL when it comes to state. 
```
curl http://localhost:3510/v1.0/state/durable-statestore/name
```
Just as a manager would like to ask the same questions of their personal assistant when it comes to remembering stuff, regardless of whether the PA writes things down on paper, memorizes them or uses a friend to retain the information, it is a fine thing for application developers to be able to use the same interaction with Dapr regardless of whether state is stored in MySQL, Redis Cache or any of the other types of state store that Dapr supports. In fact, an application developer does not need to know how and where the state will be stored and this can be changed at deployment time as the application administrator sees fit.

## 4. Telemetry

Dapr sidecars keep track of all interactions they handle. It is something like an audit trail that can be used for various purposes:
* explore dependencies
* analyze performance chararteristics and bottle necks
* report on usage of resources and applications
* solve problems

The default installation of Dapr comes with Zipkin, running in its own container. All telemetry is published by the sidecars to the Zipkin endpoint. 

Open [localhost:9411/](http://localhost:9411/) in your browser to bring up the Zipkin user interface. Click on the *Run Query* button (in the *Find a Trace* tab). A list is presented of the traces collected by Zipkin. Among them should be calls to *myapp* and *myotherapp*. Click on the *Show* button for one of the traces to see what additional information Zipkin offers. You will find details about the HTTP request and response. Click

You will find that as we add real applications that receive requests through their sidecars and that leverage Dapr's services, the value of the telemetry quickly increases. To get this information without any additional effort in either development or configuration is a big boon. 

# Dapr.io and Node applications
Dapr can be used with many technologies and from all programming languages that can make HTTP or gRPC requests. For some languages (Go, .NET, Java, Node/JavaScript, PHP and Python with support for C/C++ and Rust in active development) SDKs for Dapr have been provided that make working with Dapr easier still.

Below you will find some exercises of working with the Dapr Node SDK. Of course, similar things can easily be done from one of the other programming languages. More information on all available [SDKs in the Dapr Docs](https://docs.dapr.io/developing-applications/sdks/).

## Node Runtime Environment

Downloads of the Node runtime for various Operating Systems are available on this [download page ](https://nodejs.org/en/download/). The installation of the Node runtime will include the installation of the *npm* package manager - that will be needed for some of the demo applications.

If you are new to Node, you may want to read a little introduction on Node and its history: [Introduction to Node](https://nodejs.dev/introduction-to-nodejs). 

When the installation is done, verify its success by running

`node -v`

on the commandline. This should run successfully and return the version label for the installed Node version.

Also run:

`npm -v`

on the commandline to verify whether *npm* is installed successully. This should return the version label for the installed version of *npm*. NPM is the Node Package Manager - the component that takes care of downloading and installing packages that provide modules with reusable functionality. Check the website [npmjs.com](https://www.npmjs.com/) to explore the wealth of the more than 1 Million packages available on NPM.

## Node and Dapr
Let us find out how the Dapr sidecar and the services it can deliver can be engaged from Node applications. The code discussed in this section are in the directory *hello-world-dapr*.

The Dapr Node SDK has been installed using NPM with this statement (you do not have to execute this statement, although it does not hurt when you do):
```
npm i dapr-client --save
```
Check package.json. It includes the dependency of the Node application on the dapr-client module. The modules themselves are downloaded into the *node-modules* directory, when you execute this command (this one must be executed because otherwise the code will not work):

```
npm install
```
which reads the dependencies in package-lock.json or package.json and downloads and installs all direct (and indirect) dependencies.

Check out the *app.js* file. It contains a small application that handles HTTP requests: it stores the name passed in the request and stores it as state (in a Dapr state store). It keeps track of the number of occurrences of each name and reports in the HTTP response how many times a name has been mentioned.

The Node application is started through Dapr and as such gets a Dapr sidecar that handles the state related activities. The Node application uses the Node SDK to communicate with the Sidecar - instead of and much more convenient and elegant then explicit HTTP or gRPC interactions.  

This diagram gives an overview of the application.
![](images/app-dapr-state.png)

Function *retrieveIncrementSave* is where the real action is when it comes to leveraging the Dapr sidecar. It gets and saves state - without knowing any of the details of the state store (which for now happens to Redis Cache, but could be changed to MySQL or any type of store without any impact on the Node application). The definition of *client* is a crucial linking pin: the client connects the Node application to the Dapr side car.

The application does one other thing of interest: it reads from the state store the value stored under key *instance-sequence-number*. It increases that number (or sets it to 1 if it does not yet occur) and uses it for its own identfication. Multiple instances if this application can run - at the same time or at different points in time - and each will have their identification.

Run the application using these commands; Dapr will know the application as *nodeapp*:

```
export DAPR_HTTP_PORT=3510
export APP_PORT=3110
dapr run --app-id nodeapp  --app-port $APP_PORT --dapr-http-port $DAPR_HTTP_PORT node app.js
```
You will find that the logging from the Sidecar and the Node application appear in the same terminal window. The logging shows the identification number assigned to the currently running instance. It will probably be *one*. If you stop the application and start it again, it should be incremented by one.

Make a request to the application - you will need a second terminal window for this - a plain HTTP request directly to the application:
```
curl http://localhost:3110/?name=Joseph
```
You will get a response that indicates how often this name has occurred before. Make the same request again and find that the instance count has increased.

A different way to make the request is not directly to the Node application and the port it is listening on, but instead to the Dapr sidecar - the application's personal assistant. The sidecar can apply authorization on the request, register telemetry and perform load balancing when multiple instances of the application would be running.

The request through the sidecar is standardized into a somewhat elaborate URL:
```
curl localhost:3510/v1.0/invoke/nodeapp/method/?name=Joseph
```
The first part - localhost:3510 - refers to the Dapr sidecar and the HTTP port on which it is listening. The next segment - /v1.0/invoke - identifies the Dapr API we want to access. Subsequently we inform this API through /nodeapp that we want to interact with the application that Dapr knows as *nodeapp* and we want to pass the URL query parameter *name* with *Joseph* as its value.  

Stop the Node application and its Dapr sidecar. Ctrl+C in the terminal window where you started the application should do the trick. Then start the application again. Make the same curl call as before:
```
curl http://localhost:3110/?name=Joseph
```

This should convince you that the state written by the application survives the application. As long as the container with the Redis Cache is running, the state will be available across multiple application restarts and even multiple application instances.

### A second application

We will now add a second application to the mix. It is defined in the file *front-app.js*. This application also handles HTTP requests with a name in it. To be honest: it a very flimsy front end that has the *nodeapp* do the real work - such as name counting and state managing. The *frontapp* invokes *nodeapp*. 

Note: normally, frontapp would have its sidecar make the call to the nodeapp Dapr-application's sidecar without needing to know where it runs. However in some environments this does not seem to work as a result of incorrect name resolution. To work around this issue, we have *frontapp* call directly to the *nodeapp*'s sidecar instead of its own. This somewhat reduces the elegance of using the sidecars and is an unfortunate and tenmporary workaround.

This diagram visualizes the situation with the two applications and their sidecars.

![](images/front-app-nodeapp-statestore-sidecars.png)

Start the *frontapp* using these commands:
```
export NODE_APP_DAPR_PORT=3510
export APP_PORT=3220
export DAPR_HTTP_PORT=3620
dapr run --app-id frontapp  --app-port $APP_PORT --dapr-http-port $DAPR_HTTP_PORT node front-app.js
```
Then make a call to the *frontapp* application using curl:
```
curl localhost:3220/?name=Johnny
```
This next call to the *frontapp* through its sidecar should work - but it does not in all environments; if it fails for you, just ignore it. It probably has to do with name resolution on your system and is not important for this handson lab.
```
curl localhost:3620/v1.0/invoke/frontapp/method/greet?name=Klaas
```
Application *frontapp* has registered with Dapr and should be known to *nodeapp*'s Dapr sidecar, so this call will work - invoking *frontapp* via this sidecar for *nodeapp*:
```
curl localhost:3510/v1.0/invoke/frontapp/method/greet?name=Klaas

curl localhost:3510/v1.0/invoke/nodeapp/method/?name=Joseph
```
You should see the name occurrence increase with each call.

Now kill *nodeapp*.

Try:
```
curl localhost:3220/?name=Johnny
```
An exception is reported (because front-could not reach nodeapp). 

Restart *nodeapp*. The application instance number is increased compared to before when you make these calls - into *frontapp* (and indirectly to *nodeapp*) and directly to *nodeapp*: 

```
curl localhost:3510/v1.0/invoke/frontapp/method/greet?name=Klaas

curl localhost:3510/v1.0/invoke/nodeapp/method/?name=Joseph
```
Note that the greeting # keeps increasing: the name and the number times it has occurred is stored as state and persists across application restarts.

However, it is not ideal that frontapp depends on nodeapp in this way, and has to report an exception when nodeapp is not available.

We will make some changes:
* *frontapp* will publish a message to a pub/sub component (in Dapr, this is by default implemented on Redis)
* *nodeapp* will consume messages from the pub/sub component and will write the name to the state store and increase the occurrence count
* *frontapp* will no longer get information from *nodeapp*; it will read directly from the state store; however: it will not write to the state store, that is still the responsibility and prerogative only of *nodeapp*. 

Stop all running applications before starting the next section.

## Node and Dapr - Pub/Sub for Asynchronous Communications

Focus now on folder *hello-world-async-dapr*. It contains the app.js and front-app.js files that we have seen before - but they have been changed to handle asynchronous communications via the built in Pub/Sub support in Dapr based in this case on the out of the box Redis based message broker.

Run 
```
npm install
```
to have the required npm modules loaded to the *node-modules* directory.

Check file *~/.dapr/components/pubsub.yaml* to see how the default Pub/Sub component is configured. It gives a fairly good idea about how other brokers could be configured with Dapr, brokers such as RabbitMQ or Apache Kafka.
```
cat ~/.dapr/components/pubsub.yaml
```
The name of the component is *pubsub* and its type is *pubsub.redis*. Daprized applications will only mention the name (*pubsub*) when they want to publish or consume messages, not refer to *redis* in any way. They do not know about the *redis* subtype and if it changes (when for example a Pulsar or Hazelcast message broker is introduced), they are not impacted.

Inspect the file *consumer.js* that contains a Dapr-based message consumer application. This application constructs a Dapr Server - an object that received requests from the Dapr Sidecar. Before, we saw the Dapr Client, that is used for sending instructions to the Sidecar.

Using this DaprServer, a subscription is created for messages on topic *orders* on pubsub component *pubsub*. This subscription is provided an anonymous and asynchronous handler function that will be invoked for every message the Sidecar retrieves from the message topic. 

Run the simple sample message consuming application *order-processor*:
```
export APP_PORT=6002
export DAPR_HTTP_PORT=3602
dapr run --app-id order-processor --app-port $APP_PORT --dapr-http-port $DAPR_HTTP_PORT --dapr-grpc-port 60002 node consumer.js
```

Check the logging to find that the application is listening on HTTP port 6002 to receive any messages that the Dapr sidecar (the personal assistant to the application) may pick up based on the topic subscription.

This diagram visualizes the current situation:
![](images/consumer-subscribed-to-topic.png)

To publish a message to the *orders* topic in the default *pubsub* component, run this CLI command:
```
dapr publish --publish-app-id order-processor --pubsub pubsub --topic orders --data '{"orderId": "100"}' 
```
This tells Dapr to publish a message on behalf of an application with id *order-processor* (which is the application id of the only Dapr sidecar currently running) to the pubsub component called *pubsub* and a topic called *orders*. 

Check in the logging from the consumer application if the message has been handed over by the Dapr sidecar to the application (after consuming it from the topic on the pubsub component).

### Publishing from Node

The publisher application *orderprocessing* is a simple Node application that sends random messages to the *orders* topic on *pubsub*. Check the file *publisher.js*.  It creates a Dapr client - the connection from Node application to the Sidecar - and uses the *pubsub.publish* method on the client to publish messages to the specified TOPIC on the indicated PUBSUB component. Through the Dapr component definitions (yaml files), multiple pubsub components (backed by the same or by different providers such as Redis, RabbitMQ, Hazelcast) can be defined, each with their own name. The default components file contains the *pubsub* component, backed by Redis Cache.

Run the application with the following statement, and check if the messages it produces reach the consumer:

```
export APP_PORT=6001
export DAPR_HTTP_PORT=3601
dapr run --app-id orderprocessing --app-port $APP_PORT --dapr-http-port $DAPR_HTTP_PORT node publisher.js 
```
The publisher application is started and publishes all it has to say - to its Dapr Sidecar. This loyal assistant publishes the messages onwards, to what we know is the Redis Pub/Sub implementation.

This diagram puts it into a picture:
![](images/pub-and-sub-from-node.png)

These messages are consumed by the *consumer* app's Sidecar because of its subscription on the *orders* topic. For each message, a call is made to the handler function. 

Check the logging in the terminal window where the *consumer* app is running. You should see log entries for the messages received. Note that the messages in the log on the receiving end are not in the exact same order as they were sent in. They are delivered in the original order and each is processed in its own instance of the handler function. Since the messages in this case arrive almost at the same time and the processing times for the messages can vary slightly, the order of the log messages is not determined. 

Stop the consumer application. 

Run the publisher application again. Messages are produced. And they are clearly not received at this point because the consumer is not available for consuming them. Are these messages now lost? Has communication broken down?

Start the consumer application once more to find out:
```
export APP_PORT=6002
export DAPR_HTTP_PORT=3602
dapr run --app-id order-processor --app-port $APP_PORT --dapr-http-port $DAPR_HTTP_PORT node consumer.js
```

You should see that the messages published by the publisher application when the consumer was stopped are received by the consumer now that it is running again. This is a demonstration of asynchronous communication: two applications exchange messages through a middle man - the pubsub component - and have no dependency between them.  

The handshake between Dapr sidecar and pubsub component on behalf of the consumer is identified through the app-id. Messages are delivered only once to a specific consumer. When a new consumer arrives on the scene - with an app-id that has not been seen before - it will receive all messages the queue is still retaining on the topic in question.

Stop the consumer application.

Start the consumer application *with a new identity* - defined by the *app-id* parameter: 
```
dapr run --app-id new-order-processor --app-port $APP_PORT --dapr-http-port $DAPR_HTTP_PORT node consumer.js
```
and watch it receive all earlier published messages. 

## Leverage Dapr Pub/Sub between Front App and Node App
As was discussed before, we want to break the synchronous dependency in the front-app on the node-app. To achieve this, we will make these changes:
* the frontapp will publish a message to the *names* topic on the default pub/sub component 
* the nodeapp will consume messages from this *names* topic on the pub/sub component and will write the name from each message it consumes to the state store and increase the occurrence count for that name
* the frontapp will no longer get information from synchronous calls to the nodeapp; it will read directly the occurrence count for a name from the state store; however: it will not write to the state store, that is the task for nodeapp. 

Here we see a very simplistic application of the *CQRS* pattern where we segregate the responsibility for reading from a specific data set and writing data in that set.

The front-app.js file is changed compared to the earlier implementation:
* publish a message to the *names* topic on *pubsub* for every HTTP request that is processed
* retrieve the current count for the name received in an HTTP request from the state store (assume zero if the name does not yet occur) and use the name count increased by one in the HTTP response  

The Dapr client is used for both publishing the message and for retrieving state. The direct call from *front-app.js* to the (other) Node application has been removed.

Run the *frontapp* with these statements:

```
export APP_PORT=6030
export DAPR_HTTP_PORT=3630
dapr run --app-id greeter --app-port $APP_PORT --dapr-http-port $DAPR_HTTP_PORT node front-app.js 
```
Check in the logging that the application was successfully started.

Make a number of calls that will be handled by the front-app:
```
curl localhost:6030?name=Jonathan
curl localhost:6030?name=Jonathan
curl localhost:6030?name=Jonathan
```
You will notice that the number of occurrences of the name is not increasing. The reason: the *frontapp* cannot write to the state store and the application that should consume the messages from the pubsub's topic is not yet running and therefore not yet updating the state store. Here is an overview of the situation right now:
![](images/front-app-publisher.png)  

So let's run this *name-processor* using these statements:
 
```
export APP_PORT=6031
export SERVER_PORT=6032
export DAPR_HTTP_PORT=3631
dapr run --app-id name-processor --app-port $APP_PORT --dapr-http-port $DAPR_HTTP_PORT node app.js 
```
The logging for this application should show that the messages published earlier by *frontapp* are now consumed, and the statestore is updated. 

Here is the situation in a picture:
![](images/front-app-pub-and-nameprocessor-sub.png)

Note: this implementation is not entirely safe because multiple instances of the handler function, each working to process a different message, could end up in *race conditions* where one instance reads the value under a key from the state store, increases it and saves it. However, a second instance could have read the value right after or just before and do its own increment and save action. After both are done, the name occurrence count may be increased by one instead of two. For the purpose of this lab, we accept this possibility.    

Make a number of calls that will be handled by the front-app:
```
curl localhost:6030?name=Jonathan
curl localhost:6030?name=Jonathan
curl localhost:6030?name=Jonathan
```
You will notice that the number of occurrences of the name is (still) not increasing. However, when you check the logging for the name-processor, you should also see that it is triggered by an event that contains the name *Jonathan* and it is keeping correct count. So whay does the front-app not produce the same number of occurrences?

This has to do with different modes of operation of the Dapr state store. By default, every application has its own private area within the state store. Values stored by one application are not accessible to other applications. In this example, *name-processor* records the name and the occurence count. And when the front-app retrieves the entry for the name from the state store, it will not find it (in its own private area).

We can instruct Dapr to use a state store as a global, shared area that is accessible to all applications.  See [Dapr Docs on global state store](https://docs.dapr.io/developing-applications/building-blocks/state-management/howto-share-state/).

Copy the default state store component configuration to the local directory, as well as the *pubsub* component configuration: 
```
cp ~/.dapr/components/statestore.yaml .
cp ~/.dapr/components/pubsub.yaml .
```
Check the contents of the file that specifies the state store component:
```
cat statestore.yaml
```
You see how the state store component is called *statestore* and is of type *state.redis*. Now edit the file and add a child element under metadata in the spec (at the same level as redisHost):
```
  - name: keyPrefix
    value: none  # none means no prefixing. Multiple applications share state across different state stores
```
This setting instructs Dapr to treat keys used for accessing state in the state store as global keys - instead of application specific keys that are automatically prefixed with the application identifier.

Save the file.

Stop both the frontapp and the name-processor applications.

Start both applications - with the added components-path parameter. This parameter tells Dapr to initialize components as defined by all the yaml files in the indicated directory (in this case the current directory). That is why you had to copy the pubsub.yaml file as well to the current directory, even though it is not changed. If you would not, it is not found by Dapr and call attempts to publish messages to topics on *pubsub* or subscribe to such topics will fail.

![](images/pubsub-and-global-state.png)

In one terminal, start the *greeter* application:
```
export APP_PORT=6030
export DAPR_HTTP_PORT=3630
dapr run --app-id greeter --app-port $APP_PORT --dapr-http-port $DAPR_HTTP_PORT --components-path .  node front-app.js 
```
and in a second terminal run *name-processor*:
```
export APP_PORT=6031
export SERVER_PORT=6032
export DAPR_HTTP_PORT=3631
dapr run --app-id name-processor --app-port $APP_PORT --dapr-http-port $DAPR_HTTP_PORT --components-path . node app.js 
```

Again, make a number of calls that will be handled by the front-app:
```
curl localhost:6030?name=Michael
curl localhost:6030?name=Michael
curl localhost:6030?name=Michael
curl localhost:6030?name=Jonathan
```
At this point, the front-app should get the increased occurrence count from the state store, saved by the name-processor app, because now both apps work against the global shared state store. 

## Go applications and Dapr

Directory *dapr-go* contains a simple straightforward Go application. This application uses the Dapr Go SDK - to save state and retrieve state. Nothing very useful - but a working example of Go and the Dapr SDK.

Note: Go has great gRPC support and it is recommended to have the Go application interact with its Dapr sidecar over gRPC.

Run the application using this command:
```
dapr run --app-id orderprocessing --dapr-grpc-port 60006 go run OrderService.go
```
The value passed in a *dapr-grpc-port* is available in the Go application as environment variable and that value is used to initialize the Dapr Client that interacts from the application to the Sidecar. 

Check the logs that are produced when you execute this command. You will see Dapr starting, the application starting and looking for the Dapr Sidecar on GRPC port 60006 and then when finding it (alive and kicking) reporting that application and sidecar have joined forces and are ready for action. This action subsequently commences.

```
ℹ️  Starting Dapr with id orderprocessing. HTTP Port: 35875. gRPC Port: 60006
ℹ️  Checking if Dapr sidecar is listening on HTTP port 35875
....
WARN[0000] app channel not initialized, make sure -app-port is specified if pubsub subscription is required  app_id=orderprocessing instance=DESKTOP-NIQR4P9 scope=dapr.runtime type=log ver=edge
WARN[0000] failed to read from bindings: app channel not initialized   app_id=orderprocessing instance=DESKTOP-NIQR4P9 scope=dapr.runtime type=log ver=edge
INFO[0000] dapr initialized. Status: Running. Init Elapsed 9.1305ms  app_id=orderprocessing instance=DESKTOP-NIQR4P9 scope=dapr.runtime type=log ver=edge
INFO[0000] placement tables updated, version: 0          app_id=orderprocessing instance=DESKTOP-NIQR4P9 scope=dapr.runtime.actor.internal.placement type=log ver=edge
ℹ️  Checking if Dapr sidecar is listening on GRPC port 60006
ℹ️  Dapr sidecar is up and running.
ℹ️  Updating metadata for app command: go run OrderService.go
✅  You're up and running! Both Dapr and your app logs will appear here.
```

More details on the [Go Client SDK for Dapr](https://docs.dapr.io/developing-applications/sdks/go/go-client/)


## Telemetry, Traces and Dependencies
Open the URL [localhost:9411/](http://localhost:9411/) in your browser. This opens Zipkin, the telemetry collector shipped with Dapr.io. It provides insight in the traces collected from interactions between Daprized applications and via Dapr sidecars. This helps us understand which interactions have taken place, how long each leg of an end-to-end flow has lasted, where things went wrong and what the nature was of each interaction. And it also helps learn about indirect interactions.

![](images/zipkin-telemetery-collection.png)

Query Zipkin for traces. You should find traces that start at *greeter* and also include *name-processor*. You now that we have removed the dependency from *greeter* on *name-processor* by having the information flow via the pubsub component. How does Zipkin know that greeter and name-processor are connected? Of course this is based on information provided by Dapr. Every call made by Dapr Sidecars includes a special header that identifies a trace or conversation. This header is added to messages published to a pubsub component and when a Dapr sidecar consumes such a message, it reads the header value and reports to Zipkin that it has processed a message on behalf of its application and it includes the header in that report. Because Zipkin already received that header when the Dapr sidecar that published the message (on behalf of the greeter application) reported its activity, Zipkin can construct the overall picture.

When you go to the Dependencies tab in Zipkin, you will find a visual representation of the dependencies Zipkin has learned about. Granted, there are not that many now, but you can imagine how this type of insight in a complex network of microservices could add useful insights.


## Closure

To complete this lab, you can now stop dapr: `dapr stop`. You can also stop the MySQL container:

```
docker stop dapr-mysql
```


## Resources

[Dapr.io Docs - Getting Started](https://docs.dapr.io/getting-started/)

[Dapr.io Docs - MySQL State Store component](https://docs.dapr.io/reference/components-reference/supported-state-stores/setup-mysql/)

[Dapr.io Docs - State Management](https://docs.dapr.io/developing-applications/building-blocks/state-management/)

[MySQL Container Image docs](https://hub.docker.com/_/mysql/)