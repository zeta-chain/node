FROM zetanode:latest as zeta
FROM ethereum/client-go:v1.10.26 as geth



FROM alpine

RUN apk --no-cache add jq curl bash
COPY --from=zeta /usr/local/bin/zetacored /usr/local/bin/
COPY --from=zeta /usr/local/bin/zetaclientd /usr/local/bin/
COPY --from=geth /usr/local/bin/geth /usr/local/bin/
COPY --from=smoketest /go/bin/smoketest /usr/local/bin/

WORKDIR /work

COPY setup.js /work/
COPY start.sh /work/

RUN chmod +x /work/start.sh
RUN chmod +x /usr/local/bin/smoketest

