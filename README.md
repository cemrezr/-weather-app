# Weather App

Weather App is a Golang-based application that provides weather information for a specified location using various weather services. The application fetches weather data from external APIs and stores them in a PostgreSQL database.

## Technical Overview

### Endpoint: 
The application provides a single endpoint `/weather?q=\<location\>`.  For each request to this endpoint, the app fetches data from two external weather APIs, averages the temperature values.

### Request Handling:
To reduce the cost of using external weather APIs, incoming requests for the same location are grouped and handled together.  If multiple requests for the same location are received within a 5-second window, they are consolidated into a single API call, and the response is shared among all waiting requests.

### Request Grouping and Batching:
All incoming requests for the same location within a 5-second timeframe are handled together.
If a request is received after 3 seconds (midway through the 5-second window), it will be grouped with ongoing requests, resulting in a slightly longer wait time for this request.    
A maximum of 10 concurrent requests per location can be grouped. If the number exceeds this limit, the grouped requests are sent immediately without waiting for the remaining duration.


## Prerequisites

Before you begin, ensure you have the following tools installed on your system:

- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Project Setup
Build and run the application using Docker Compose:
To build and start the application, run the following command:

```


docker-compose up --build

```

This command will build the application image, create the necessary containers, and start the services in detached mode.

## Testing the Application

Access the Application:
You can access the weather application via the following endpoint:

```

http://localhost:8081//weather?q=<location>

```

Replace {location} with the desired city name (e.g., http://localhost:8081/weather/Istanbul).
Example request:

```

curl http://localhost:8080/weather?q=Istanbul
```

This should return the average temperature for the specified location.

If you want to send multiple requests simultaneously and run them in the background, you can use the & operator at the end of the curl commands as shown below:

```
curl http://localhost:8080/weather?q=Istanbul
curl http://localhost:8080/weather?q=Istanbul
curl http://localhost:8080/weather?q=Istanbul

```


## Inspect the Database
If you want to verify that the data is being stored in the database, you can log into the PostgreSQL container:

```

docker-compose exec postgres psql -U weather_user -d weather_db

```

Then, check the tables and records:

```

SELECT * FROM weather_queries;

```
