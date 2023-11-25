FROM echogosdkcompiler:latest

ENV TARGET=armv7l-linux-gnueabi

# Install some dependencies
RUN apt-get install autoconf libtool -y

# Build libelf
WORKDIR /
RUN git clone https://android.googlesource.com/platform/external/elfutils --single-branch --branch android-7.1.1_r21
WORKDIR /elfutils
# Some addtional project setup is required
RUN cp -a /elfutils/bionic-fixup/. $NDK_INCLUDE
RUN sed -i 's\-Wstack-usage=262144\\g' config/eu.am
RUN sed -i '294,320d' configure.ac
RUN sed -i '286,291d' configure.ac
RUN sed -i '87,109d' configure.ac
ENV CFLAGS="$CFLAGS -Wno-error=implicit-function-declaration -Wno-error=int-conversion -include /elfutils/bionic-fixup/AndroidFixup.h -I/elfutils/bionic-fixup/"
#Start compiling process
RUN autoreconf -fvi
RUN ./configure --host $TARGET --disable-dependency-tracking --includedir=/elfutils/bionic-fixup --oldincludedir=/elfutils/bionic-fixup
RUN cd libelf && make
# compiled: /elfutils/libelf/libelf.so

# Build ltrace
WORKDIR /
RUN git clone https://android.googlesource.com/platform/external/ltrace --single-branch --branch android-7.1.1_r21
WORKDIR /ltrace
# Additional setup
RUN sed -i 's\string rindex(string,char);\string strrchr(string,char);\g' etc/libc.so.conf
RUN sed -i 's\rindex\strrchr\g' ltrace-elf.c
# Link with our built libelf
RUN mv /elfutils/libelf/libelf.so $NDK_LIBS/libelf.so
RUN cp -a /elfutils/libelf/. $NDK_INCLUDE
# Build ltrace
ENV CFLAGS="-Wno-error=string-plus-int -Wno-error=implicit-function-declaration -Wno-error=duplicate-decl-specifier -Wno-error=deprecated-declarations"
RUN ./autogen.sh
RUN ./configure --host $TARGET
RUN make
RUN make install

# Finished compiling
RUN mkdir -p /ltrace/build
CMD ["cp", "ltrace", "build/ltrace"]