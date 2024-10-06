# Weather App

Weather App is a Golang-based application that provides weather information for a specified location using various weather services. The application fetches weather data from external APIs and stores them in a PostgreSQL database.


## Prerequisites

Before you begin, ensure you have the following tools installed on your system:

- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Project Setup

Running the Application

Build and Start the Application
To build and start the application, run the following command:

bash
docker-compose up -d --build

This command will build the application image, create the necessary containers, and start the services in detached mode.

Testing the Application

Access the Application
You can access the weather application via the following endpoint:

http://localhost:8081/weather/{location}

Replace {location} with the desired city name (e.g., http://localhost:8081/weather/Istanbul).
Example request:
curl http://localhost:8081/weather/Istanbul

This should return the average temperature for the specified location.

Inspect the Database
If you want to verify that the data is being stored in the database, you can log into the PostgreSQL container:

docker-compose exec postgres psql -U weather_user -d weather_db

Then, check the tables and records:
SELECT * FROM weather_queries;
