# merc-benz-route-checker

## Intro

The application is to identify the minimum number of stations required to stop for recharging your car when travelling from a source to a destination.

## Tech stack

* Language: golang
* Http web framework: [Gin](https://github.com/gin-gonic/gin)
* Docker containerization

## Running the app (without containerization)

The below command generates a linux build under the path `./dist/benz`

`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/benz "main.go"`

After building the app, run the app with the command `./benz` under the `./dist` folder.

## Running the app using docker

### Building docker image

Use the below command to build the docker using `Dockerfile`.

`docker build . -t benz-service --build-arg mode=prod`

### Running the container

Use the below command to run the docker image that was built in the previous command.

`docker run -p 8080:8080 --name benz-service <image id>`

Once the docker container is up, the APIs can be accessed under the host address `http://localhost:8080`

## APIs

The below APIs are available in the service

* `http://localhost:8080/api/health` - health check API
* `http://localhost:8080/api/v1/compute-route` - API to compute route with minimum number of stops

## About the logic to find the minimum number of charging station

The algorithm uses greedy approach with priority queue. With the greedy approach, we charge the car only at stations that can provide maximum number of charges when compared to all other stations at that state. Below is the logical explanation of the method.
1. We find out the maximum distance the car can travel with available charge.
2. If the destination can be reached with available charge, we return an empty slice. It means that there is no need for the car to stop for recharging since the available charge is sufficient.
3. If the destination cannot be reached with available charge, the logic simulates the car to travel to maximum distance possible noting down the stations along the route in priority queue.
4. The priority queue will be in descending order respect to the charge available in station. For example, if the station and charge pair are S1:10, S2:20, S3:30, then the priority queue will return in the order S3:30, S2:20, S1:10. We always pick the next station that provides maximum charge.
5. If the charge in a station is not sufficient, we pick the next station from the priority queue. This is done till either the queue is empty or the charge becomes sufficient.
6. The station names where the car recharges are added to the returning slice.

The time complexity of this logic is O(nlog(n)). We iterate n times and greedily check if recharge is required.

The space complexity of this logic is O(n)


## Support for ELK stack

ELK stack is used to push and analyse logs. If the microservices environment supports ELK stack, then this service can be leveraged/configured to work with it. The code is available in the service but commented out.

## Support for instrumentation

Code for instrumentation is added. If the microservices environment supports instrumentation/observability tool like graphite, this microservice can be configured to it.

## Adding CI/CD jenkins pipeline

The above steps can be covered in jenkins pipeline to automate the build and deploy it into the microservice environment.

## Addition of .env files

The `app-dev.env` and `app-prod.env` files haven't been removed for reference.


MODE=prod SHIPLOGS=true GRAPHITE_URL=graphite:8125 LOGSTASH_URL=logstash:8089 docker-compose up -d 