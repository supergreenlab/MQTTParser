FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD supergreenlog /

EXPOSE 8080

ENTRYPOINT ["/supergreenlog"]
