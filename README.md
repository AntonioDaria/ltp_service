# Last Traded Price API

A simple API to fetch the last traded price of cryptocurrency pairs from the Kraken API.


```
As a user
so I can get the last traded price of cryptocurrency pairs
i can send a GET request to `http://localhost:8080/api/v1/ltp`
```

## What's been accomplished ✅

- The API, fetches the last traded price from the Kraken API, and returns it in the response for the following currency pairs:
1.BTC/USD
2.BTC/CHF
3.BTC/EUR

## Nice to haves 💫

- Add support for more cryptocurrency pairs.
- Make the api dynamic to accept a request body with a chosen currency pairs


## Getting started

To run it locally 

### Prerequisites

- Docker
- Docker Compose
- Make

### Building and Running the Application

1. **Clone the repository:**
```sh
git clone <repository-url>
cd <repository-directory>
```
   
2. **Ensure you have Swag CLI installed:**
```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

3. **Start the application:**
```sh
make swag 

make start
```

## Interacting with the API
```
You can test the API using Swagger UI or Postman:

Using Swagger UI
Visit http://localhost:8080/swagger/index.html to access the Swagger UI and interact with the API.

Using Postman
- URL: http://localhost:8080/api/v1/ltp
- Method: GET
```

### Technologies used
- Golang: The primary language for the API.
- Fiber: Web framework for Golang.
- Gomock: Mocking framework for unit tests.
- net/http/httptest: Package for HTTP testing.
- stretchr/testify: Assertion library for testing.


### Testing
All current tests are passing. Tests use Gomock and stretchr/testify.
This will run all the tests in the project.


3. **To Run the tests:**
```sh
go test ./...
```

4. **Stop the application:**
```sh
make stop
```

Conclusion
By following these instructions, you should be able to build, run, and interact with the Last Traded Price API. If you encounter any issues or have any questions, feel free to contact antoniodaria@msn.com.