FROM eu.gcr.io/antha-images/golang:1.12.4-build AS build
ARG COMMIT_SHA
ARG NETRC
RUN printf "%s\n" "$NETRC" > /root/.netrc
RUN mkdir /antha
WORKDIR /antha
RUN set -ex && go mod init antha && go mod edit "-require=github.com/antha-lang/antha@$COMMIT_SHA" && go mod download
# Do these builds to pre-warm the build cache. This makes a HUGE difference to performance in the cloud
RUN set -ex && go build github.com/antha-lang/antha/... github.com/Synthace/antha-runner/... github.com/Synthace/instruction-plugins/... && go install github.com/antha-lang/antha/cmd/...
RUN set -ex && go test -c github.com/antha-lang/antha/cmd/elements
COPY scripts/*.sh /antha/

FROM eu.gcr.io/antha-images/golang:1.12.4-build AS tests
COPY --from=build /root/.netrc /root/.cache /root/
COPY --from=build /go /go
COPY --from=build /antha /antha
WORKDIR /antha
RUN ./antha-test.sh

FROM eu.gcr.io/antha-images/golang:1.12.4-build AS cloud
## This target produces an image that is used both for gitlab elements CI, and also workflow execution in the cloud
COPY --from=tests /root/.cache /root/
COPY --from=tests /go /go
COPY --from=build /antha /antha
WORKDIR /antha
# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
