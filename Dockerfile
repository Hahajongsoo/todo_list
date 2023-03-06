FROM ubuntu:latest
RUN apt update && apt install -y ca-certificates
WORKDIR /app

COPY todoWeb /app/
COPY public /app/public

EXPOSE 3000

ENTRYPOINT [ "./todoWeb" ]