FROM ubuntu:latest
MAINTAINER XinWang

RUN apt update -y
RUN apt install make imagemagick -y
WORKDIR /
#RUN wget https://golang.google.cn/dl/go1.17.8.linux-armv6l.tar.gz
#RUN LS /
#RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.8.linux-arm6l.tar.gz
ADD . /dkmission
#ENV PATH=$PATH:/usr/local/go/bin:
#RUN go env -w GO111MODULE=on
#RUN go env -w GOPROXY=https://goproxy.cn,direct
ENV LD_LIBRARY_PATH=/dkmission/processor 
WORKDIR /dkmission 
RUN make dockerPrepare -j8


