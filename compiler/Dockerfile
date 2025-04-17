FROM golang:1.24.0

LABEL org.opencontainers.image.source = "https://github.com/Binozo/EchoGo"

ENV NDK_VERSION=21.4.7075529
ENV API=24
ENV GRADLE_VERSION=6.3
ENV CMD_TOOLS_VERSION=6200805

ENV GOOS=android
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=1

# Base setup
WORKDIR /
RUN apt-get update && \
    apt-get dist-upgrade -y && \
    apt-get install -y make git default-jdk unzip wget

# Install Android NDK
RUN wget https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip -P /tmp && \
    unzip -d /opt/gradle /tmp/gradle-${GRADLE_VERSION}-bin.zip && \
    rm /tmp/gradle-${GRADLE_VERSION}-bin.zip

RUN mkdir -p /opt/android/cmdline-tools
RUN wget "https://dl.google.com/android/repository/commandlinetools-linux-${CMD_TOOLS_VERSION}_latest.zip" -P /tmp && \
    unzip -d /opt/android/cmdline-tools /tmp/commandlinetools-linux-${CMD_TOOLS_VERSION}_latest.zip && \
    rm /tmp/commandlinetools-linux-${CMD_TOOLS_VERSION}_latest.zip

# Accept licenses & install tools using the sdkmanager
RUN yes Y | /opt/android/cmdline-tools/tools/bin/sdkmanager --install "build-tools;29.0.3" "platforms;android-${API}" "platform-tools" "ndk;$NDK_VERSION" && \
    yes Y | /opt/android/cmdline-tools/tools/bin/sdkmanager --licenses

# Set some environment variables
ENV GRADLE_HOME=/opt/gradle/gradle-${GRADLE_VERSION}
ENV ANDROID_HOME=/opt/android
ENV ANDROID_NDK_HOME=/opt/android/ndk/$NDK_VERSION
ENV PATH="$PATH:/opt/gradle/gradle-${GRADLE_VERSION}/bin:/opt/gradlew:/opt/android/emulator:$ANDROID_HOME/cmdline-tools/tools/bin:/opt/android/platform-tools:/opt/android/ndk/$NDK_VERSION"
ENV LD_LIBRARY_PATH="$ANDROID_HOME/emulator/lib64:$ANDROID_HOME/emulator/lib64/qt/lib"
ENV TOOLCHAIN=$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/linux-x86_64
ENV AR=$TOOLCHAIN/bin/llvm-ar
ENV AS=$TOOLCHAIN/bin/armv7a-linux-androideabi${API}-clang
ENV LD=$TOOLCHAIN/bin/arm-linux-androideabi-ld
ENV CC=$TOOLCHAIN/bin/armv7a-linux-androideabi${API}-clang
ENV CXX=$CC++
ENV RANLIB=$TOOLCHAIN/bin/llvm-ranlib
ENV STRIP=$TOOLCHAIN/bin/llvm-strip
ENV TARGET=armv7a-linux-androideabi
ENV NDK_INCLUDE=$TOOLCHAIN/sysroot/usr/include
ENV NDK_CXX_INCLUDE=$TOOLCHAIN/sysroot/usr/include/c++/v1
ENV NDK_LIBS=$TOOLCHAIN/sysroot/usr/lib/arm-linux-androideabi/$API
#RUN rm $NDK_LIBS/../libc++_shared.so

# Download includes
WORKDIR /
RUN git clone https://android.googlesource.com/platform/frameworks/native --single-branch --branch android-7.1.2_r23 && \
    cp -r /native/include/* $NDK_INCLUDE && \
    rm -r native && \
    git clone https://android.googlesource.com/platform/frameworks/native --single-branch --branch android-4.2.2_r1 && \
    mkdir -p $NDK_INCLUDE/utils && \
    cp -r /native/include/utils/* $NDK_INCLUDE/utils && \
    git clone https://android.googlesource.com/platform/system/core --single-branch --branch android-7.1.2_r39 && \
    mkdir -p $NDK_INCLUDE/cutils && \
    mkdir -p $NDK_INCLUDE/log && \
    mkdir -p $NDK_INCLUDE/system && \
    cp -r /core/include/cutils/* $NDK_INCLUDE/cutils && \
    cp -r /core/include/log/* $NDK_INCLUDE/log && \
    cp -r /core/include/system/* $NDK_INCLUDE/system && \
    git clone https://android.googlesource.com/platform/libnativehelper --single-branch --branch android-7.1.2_r39 && \
    mkdir -p $NDK_INCLUDE/nativehelper && \
    cp -r /libnativehelper/include/nativehelper/* $NDK_INCLUDE/nativehelper && \
    wget -O $NDK_INCLUDE/linux/binder.h https://android.googlesource.com/platform/prebuilts/ndk/+/refs/tags/android-7.1.2_r39/current/platforms/android-24/arch-arm/usr/include/linux/binder.h?format=TEXT

#RUN git clone https://android.googlesource.com/platform/prebuilts/ndk --branch android-7.1.2_r39 && \
#    mkdir -p $NDK_INCLUDE/linux && \
#    cp -r /ndk/current/platforms/android-24/arch-arm/usr/include/linux/* $NDK_INCLUDE/linux

WORKDIR /

# In order to correctly link we need to build tinyalsa ourself
RUN git clone https://github.com/tinyalsa/tinyalsa.git
WORKDIR /tinyalsa
# Prepare our header files
RUN make && \
    make install && \
    ldconfig

# Now we need to build the shared library (libtinyalsa.so)
WORKDIR /tinyalsa/src
# In order to link it correctly we modify the Makefile a bit
# (We need a library with any version number after .so)
RUN sed -i 's\-shared -Wl,-soname,\-shared -soname=\g' Makefile && \
    sed -i 's\libtinyalsa.so.$(LIBVERSION_MAJOR)\libtinyalsa.so\g' Makefile && \
    sed -i 's\libtinyalsa.so.$(LIBVERSION)\libtinyalsa.so\g' Makefile && \
    sed -i 's\ln -sf\#ln -sf\g' Makefile && \
    make clean

# Now we compile it using the Android NDK
RUN make CC=$CC LD=$LD AR=$AR AS=$AS libtinyalsa.a libtinyalsa.so

# Lets arrange some libraries and header files
RUN cp -r /usr/local/include/tinyalsa/ /usr/include/ && \
    cp -r /usr/include/tinyalsa/ $NDK_INCLUDE && \
    cp /tinyalsa/src/libtinyalsa.so $NDK_LIBS/libtinyalsa.so

# Finished preparing our compiler

WORKDIR /sdk
CMD [ "go", "build", "-tags='server'", "-o", "build/server", "cmd/server.go" ]