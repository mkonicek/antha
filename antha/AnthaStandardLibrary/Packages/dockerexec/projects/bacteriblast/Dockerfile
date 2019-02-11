FROM ubuntu:16.04

RUN apt-get update -y
RUN apt-get install -y curl
RUN apt-get install -y ncbi-blast+

RUN useradd -m biouser

WORKDIR /home/biouser

RUN chown -R biouser:biouser .

USER biouser

WORKDIR /home/biouser

RUN mkdir mapped
COPY example.fa ./

RUN curl ftp://ftp.ensemblgenomes.org/pub/release-42/bacteria//fasta/bacteria_32_collection/escherichia_coli_179100/cdna/Escherichia_coli_179100.ASM35690v1.cdna.all.fa.gz > Escherichia_coli_179100.ASM35690v1.cdna.all.gz

RUN curl ftp://ftp.ensemblgenomes.org/pub/release-42/bacteria//fasta/bacteria_91_collection/escherichia_coli_1303/cdna/Escherichia_coli_1303.ASM82998v1.cdna.all.fa.gz > Escherichia_coli_1303.ASM82998v1.cdna.all.gz

RUN gunzip *.gz

RUN makeblastdb -dbtype=nucl -in Escherichia_coli_179100.ASM35690v1.cdna.all
RUN makeblastdb -dbtype=nucl -in Escherichia_coli_1303.ASM82998v1.cdna.all

