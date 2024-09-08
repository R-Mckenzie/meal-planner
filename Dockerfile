FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main /app/cmd/api/main.go

FROM scratch
COPY --from=build /app/main ./main
EXPOSE 8000
CMD ["/main"]
