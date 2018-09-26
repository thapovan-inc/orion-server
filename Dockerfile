FROM scratch
ADD orion-server /orion-server
ADD default.toml /default.toml
EXPOSE 9071/tcp
EXPOSE 20691/tcp
CMD ['/orion-server']