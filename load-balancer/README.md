
# Load Balancer

The Load Balancer is a Go-based application that distributes incoming requests across multiple backend servers using different load balancing algorithms. This project aims to demonstrate three load balancing algorithms: Round Robin, Least Connections, and Random.

## Prerequisites

- Go programming language (version 1.16+)

## Installation

1. Clone the repository: \
`git clone https://github.com/your-username/load-balancer.git`

2. Change to the project directory:\
   `cd load-balancer`

3. Build the application:\
`go build`

4. Run the application: \
`./load-balancer`

## Usage

The load balancer service will start running on port 8080. You can send HTTP requests to the load balancer endpoints to observe the load balancing behavior.

### Endpoints

The load balancer provides the following endpoints:

- `/round-robin`: Uses the Round Robin load balancing algorithm.
- `/least-connections`: Uses the Least Connections load balancing algorithm.
- `/random`: Uses the Random load balancing algorithm.

### Example

Assuming the load balancer is running on `http://localhost:8080`, you can make requests to the load balancer endpoints using tools like cURL or web browsers:

1. Round Robin load balancing:\
`curl http://localhost:8080/round-robin`

2. Least Connections load balancing:\
`curl http://localhost:8080/least-connections`

3. Random load balancing:\
`curl http://localhost:8080/random`

### TODO

1. Add unit test cases.
2. Add logs, validations and error handling.



