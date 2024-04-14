FROM golang:1.22.2-alpine AS BuildStage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth ./cmd/auth

FROM alpine:latest 

COPY --from=BuildStage /app/auth /app/auth

COPY --from=BuildStage app/config/config.yaml /config/config.yaml

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT [ "/app/auth" ]



