FROM golang:latest
RUN mkdir /app
COPY /app /app
COPY go.mod /app
COPY go.sum /app
WORKDIR /app
RUN ls ./key_distributor
RUN go mod download
RUN go build -o output .
EXPOSE 13800
CMD ["/app/output"]