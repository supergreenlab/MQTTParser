FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD supergreenlog /

CMD ["/supergreenlog"]
