# GoReleaser pre-builds the binary; this stage just packages it.
FROM alpine:3.20
RUN adduser -D -h /home/envii envii
COPY envii /usr/local/bin/envii
USER envii
WORKDIR /home/envii
ENTRYPOINT ["envii"]
