FROM golang:1.22 AS build
COPY ./ /app/
WORKDIR /app/
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/nine-shop

FROM alpine:3.12
ENV TZ=Asia/Bangkok
COPY --from=build /bin/nine-shop /



EXPOSE 2525
CMD ["/nine-shop"]
