FROM golang:1.19-alpine AS build

WORKDIR /project

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/app ./cmd/app

FROM golang:1.19-alpine 

WORKDIR /project

COPY --from=build /project/bin .
COPY --from=build /project/database_init.sql .

EXPOSE 3000

ENTRYPOINT ["/project/app"]