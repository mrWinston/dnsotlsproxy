FROM golang:1.11 AS dev

WORKDIR /code
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app .

ENTRYPOINT ["go", "run"]

FROM scratch
COPY --from=dev /etc/ssl/certs /etc/ssl/certs
COPY --from=dev /app /app 
CMD ["/app"]
