FROM golang:1.16-alpine As builder

RUN mkdir /app_dir
COPY . /app_dir
WORKDIR /app_dir
RUN CGO_ENABLED=0 go build -o /app .

FROM scratch
COPY --from=builder /app /app
CMD ["./app"]
