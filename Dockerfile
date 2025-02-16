FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk update && apk add --no-cache git
WORKDIR /avito-shop-service
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /avito-shop-service/cmd/merch
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o avito-shop-service/main .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder avito-shop-service/. .
WORKDIR /avito-shop-service/cmd/merch
COPY --from=builder avito-shop-service/cmd/merch .
EXPOSE 8081
CMD ["avito-shop-service/main"]