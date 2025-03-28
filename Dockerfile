FROM golang:1.22 AS build
WORKDIR /github.com/lzchong/receipt-processor
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/receipt_processor cmd/api/main.go

FROM scratch
COPY --from=build /bin/receipt_processor /bin/receipt_processor
EXPOSE 8080
CMD ["/bin/receipt_processor"]
