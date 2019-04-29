FROM eu.gcr.io/antha-images/golang:1.12.4-build

ARG COMMIT_SHA
ADD .netrc .coveralls_token /
RUN mv /.netrc /.coveralls_token $HOME/
RUN mkdir /antha
WORKDIR /antha
RUN set -ex && go get -u github.com/mattn/goveralls
RUN set -ex && go mod init antha && go get github.com/antha-lang/antha@$COMMIT_SHA && go mod download
RUN set -ex && go install github.com/antha-lang/antha/cmd/...
RUN set -ex && go test -c github.com/antha-lang/antha/cmd/elements
COPY scripts/. /antha/.
RUN ./antha-test.sh
RUN rm $HOME/.netrc $HOME/.coveralls_token

# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
