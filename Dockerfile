FROM scratch
ADD orion-server /orion-server
ADD default.toml /default.toml
EXPOSE 9071,20691
ENTRYPOINT ['/orion-server']