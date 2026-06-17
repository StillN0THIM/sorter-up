FROM golang:1.22-alpine AS builder

WORKDIER /build

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=;inux go build -o inference-api ./cmd/api


FROM alpine:3.19

WORKEDIR /app

RUN apk add --no-cache ca-certificates 

COPY --from=builder /build/inference-api .

EXPOSE 8080
CMD ["./inference-api"]