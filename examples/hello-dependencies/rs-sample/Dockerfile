FROM rust:slim-bullseye AS builder

WORKDIR /app
COPY . /app/

RUN cargo build --release

FROM debian:bullseye-slim
COPY --from=builder /app/target/release/rs-sample /app/

RUN useradd -m appuser
USER appuser
WORKDIR /app
CMD ["./rs-sample"]