FROM golang:1.25-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o ./api ./main.go

FROM alpine:latest AS final
COPY --from=build /app/api .
EXPOSE 8080
CMD ["./api"]