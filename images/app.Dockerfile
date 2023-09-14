FROM golang:1.21 AS base

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download


FROM base AS telegram-callback

COPY telegram-callback/*.go ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /telegram-callback

EXPOSE 80

CMD ["/telegram-callback"]


FROM base AS telegram-message-processor

COPY telegram-message-processor/*.go ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /telegram-message-processor

CMD ["/telegram-message-processor"]


FROM base AS telegram-sender

COPY telegram-sender/*.go ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /telegram-sender

CMD ["/telegram-sender"]