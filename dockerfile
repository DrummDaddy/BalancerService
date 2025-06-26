FROM golang:1.23-alpine 

WORKDIR /app 

COPY go.mod go.sum ./ 

RUN go mod download 

COPY . .

RUN go build -o balancer-service main.go 

EXPOSE 443 

CMD ["./balancer-service"] 
