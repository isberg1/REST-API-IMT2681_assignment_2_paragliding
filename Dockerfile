FROM golang:1.11 as base

WORKDIR /app/

COPY . .

RUN go mod download
 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch

WORKDIR /app/

COPY --from=base /app/app .

EXPOSE 8080

CMD ["./app"]