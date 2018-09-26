FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD supergreenlog /

ENTRYPOINT ["/supergreenlog"]
