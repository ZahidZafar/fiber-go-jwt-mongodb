# Go Fiber with JWT and MongoDB Integration

This repository showcases a simple Go Fiber application integrated with JWT (JSON Web Tokens) for authentication and MongoDB for data storage. The application is containerized using Docker for easy deployment.

## Features

- User authentication using JWT.
- CRUD (Create, Read, Update) operations with a MongoDB database.
- Docker setup for easy deployment and testing.

## Technologies Used

- [Go](https://golang.org/) - The programming language used.
- [Fiber](https://github.com/gofiber/fiber) - A web framework for Go.
- [MongoDB](https://www.mongodb.com/) - A popular NoSQL database.
- [JWT](https://jwt.io/) - JSON Web Tokens for secure communication.
- [Docker](https://www.docker.com/) - Containerization for easy deployment and testing.

## Getting Started

### Prerequisites

- Install Go by following the [official Go installation guide](https://golang.org/doc/install).
- Install Docker by following the [official Docker installation guide](https://docs.docker.com/get-docker/).
- Download the MongoDB Docker image:  

```bash
   docker pull mongo
```

### Setup

1. Clone the repository:

```bash
   git clone https://github.com/ZahidZafar/fiber-go-jwt-mongodb.git
   cd your_repository
```

2. Configure environment variables:
   Create a `app.env` file in the project root and set the following variables:

```bash
   MONGO_DB_URL=mongodb://your_mongo_url
   TOKEN_SECRET_KEY=your_jwt_secret
   ACCESS_TOKEN_DURATION=25m
   TEMP_TOKEN_DURATION=1m
   REFRESH_TOKEN_DURATION=24h
```

3. Build go-app image:

```bash
   docker build -t go-app:0.1 .
```

4. Run the Docker container:

```bash
   docker-compose -f docker-compose.yaml up
```

The application will be accessible at `http://localhost:3000`.

## Contribution

Contributions are welcome! Feel free to open an issue or submit a pull request.
