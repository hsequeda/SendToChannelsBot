FROM golang:1.16-alpine As builder

RUN mkdir /app_dir
COPY . /app_dir
WORKDIR /app_dir
RUN go build -ldflags "-linkmode external -extldflags -static" -a -o /app .

FROM scratch
COPY --from=builder /app /app
CMD ["bash"]
