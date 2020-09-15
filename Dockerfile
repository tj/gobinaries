FROM golang:1.15
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -o /server cmd/gobinaries-api/main.go
CMD [ "/server" ]