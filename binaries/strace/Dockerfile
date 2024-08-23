FROM ghcr.io/binozo/echogosdkcompiler:latest

ENV TARGET=armv7l-linux-gnueabi

# Install some dependencies
RUN apt-get update
RUN apt-get install autoconf libtool build-essential flex bison -y

# Build strace
WORKDIR /
RUN git clone https://android.googlesource.com/platform/external/strace --single-branch --branch android-7.1.1_r21
WORKDIR /strace
RUN sed -i 's\# include "printstat.h"\\g' file.c
RUN sed -i 's\#include "printstat.h"\\g' file.c
RUN sed -i 's\do_printstat\;//do_printstat\g' file.c
RUN ./bootstrap
RUN ./configure --host $TARGET
RUN make
#RUN make install

# Finished compiling
RUN mkdir -p /strace/build
CMD ["cp", "strace", "build/strace"]