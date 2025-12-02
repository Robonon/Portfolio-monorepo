FROM golang:1.25-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o ./api ./main.go

FROM alpine:latest AS final
WORKDIR /app
# RUN addgroup -S -g 10001 appgroup && adduser -S -u 10001 -G appgroup appuser
COPY --from=build /app/api .
COPY ./generator/schemas.json ./schemas.json
RUN apk add --no-cache git
RUN mkdir -p /app/output
# RUN chown -R appuser:appgroup /app /app/output
# USER appuser
# WORKDIR /app/output
# RUN git init && git config user.name "Auto Commit" && git config user.email "auto@example.com"
WORKDIR /app
EXPOSE 8080
CMD ["./api"]