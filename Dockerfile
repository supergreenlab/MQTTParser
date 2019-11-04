FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD bin/supergreenlog /

EXPOSE 8080

ENTRYPOINT ["/supergreenlog"]
