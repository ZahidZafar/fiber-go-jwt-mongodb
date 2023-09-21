FROM golang:latest

WORKDIR /app


# ENV MONGO_DB_USERNAME=admin \
#     MONGO_DB_USERNAME=password

COPY go.mod go.sum ./
RUN go mod download

COPY  . .


RUN go build -o main

EXPOSE 3000


CMD [ "./main" ]
