FROM ubuntu:16.04

RUN apt-get update -y
RUN apt-get install -y build-essential g++ cmake git-all

RUN useradd -m biouser

WORKDIR /home/biouser

RUN git clone https://github.com/primer3-org/primer3.git primer3 

WORKDIR /home/biouser/primer3/src

RUN make
#RUN make test

ENV PATH=$PATH:/home/biouser/primer3/src