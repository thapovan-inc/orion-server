FROM scratch
ADD orion-server /orion-server
ADD default.toml /default.toml
ENTRYPOINT ['/orion-server']