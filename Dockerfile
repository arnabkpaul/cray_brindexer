FROM dtr.dev.cray.com/baseos/centos:7-build25 as buildenv

WORKDIR /app

ARG GOIMAGE=go1.12.7.linux-amd64.tar.gz
RUN /bin/ls -aCF && \
 yum install -y \
   wget \
   tar \
   gzip \
   git \
   gcc \
   binutils \
   e2fsprogs \
   openssl-devel \
   bzip2-devel \
   libffi-devel \
   make \
   libyaml-devel \
   && \
 cd /usr/src && \
 wget https://storage.googleapis.com/golang/$GOIMAGE &&\
 tar -C /usr/local -xzf $GOIMAGE &&\
 rm -f $GOIMAGE

RUN cd /app && \
  wget -cr -nd --no-parent -A 'lustre-client-2.12*.rpm' \
      http://steve-0.dev.cray.com/storage/lustre_builds/regular-build/cray-2.12-int/client/latest/ && \
  mkdir lustre && \
  mv lustre-client-2.12*.rpm lustre && \
  rpm -ivh --nodeps lustre-client-2.12*.rpm

FROM buildenv as builder

WORKDIR /app

ADD ./ /app/
ENV PATH=$PATH:/usr/local/go/bin
RUN cd ./go && ./prepare && ./build

FROM dtr.dev.cray.com/baseos/centos:7-build25 as production

WORKDIR /app
RUN mkdir -p /app/go/bin
COPY --from=builder /app/go/bin /app/go/bin
COPY --from=builder /app/start.sh /app/start.sh
COPY --from=buildenv /app/lustre /app
RUN rpm -ivh --nodeps lustre-client-2.12*.rpm && \
    rm -rf lustre* &&\
    yum install -y libyaml-devel

CMD ./start.sh
