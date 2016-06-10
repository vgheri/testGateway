FROM scratch
EXPOSE 1337
ADD testgateway /
CMD ["/testgateway"]
