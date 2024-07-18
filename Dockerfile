FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go get -d -v ./...
RUN go build -o mealplanner
EXPOSE 8080
CMD ["./mealplanner"]
