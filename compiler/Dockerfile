FROM golang:1.20

# Base setup
WORKDIR /
RUN apt-get update
RUN apt-get dist-upgrade -y
RUN apt-get install make git default-jdk unzip wget -y

# Install Android NDK
RUN wget https://services.gradle.org/distributions/gradle-6.3-bin.zip -P /tmp && unzip -d /opt/gradle /tmp/gradle-6.3-bin.zip
RUN mkdir -p /opt/android/cmdline-tools
RUN wget 'https://dl.google.com/android/repository/commandlinetools-linux-6200805_latest.zip' -P /tmp && unzip -d /opt/android/cmdline-tools /tmp/commandlinetools-linux-6200805_latest.zip
ENV NDK_VERSION=21.4.7075529
# Accept licenses & install tools using the sdkmanager
RUN yes Y | /opt/android/cmdline-tools/tools/bin/sdkmanager --install "build-tools;29.0.3" "platforms;android-24" "platform-tools" "ndk;$NDK_VERSION"
RUN yes Y | /opt/android/cmdline-tools/tools/bin/sdkmanager --licenses
# Set some environment variables
ENV GRADLE_HOME=/opt/gradle/gradle-6.3
ENV ANDROID_HOME=/opt/android
ENV ANDROID_NDK_HOME=/opt/android/ndk/$NDK_VERSION
ENV PATH="$PATH:/opt/gradle/gradle-6.3/bin:/opt/gradlew:/opt/android/emulator:$ANDROID_HOME/cmdline-tools/tools/bin:/opt/android/platform-tools:/opt/android/ndk/$NDK_VERSION"
ENV LD_LIBRARY_PATH="$ANDROID_HOME/emulator/lib64:$ANDROID_HOME/emulator/lib64/qt/lib"
ENV API=24
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
RUN rm $NDK_LIBS/../libc++_shared.so

# Clean up
RUN rm /tmp/commandlinetools-linux-6200805_latest.zip
RUN rm /tmp/gradle-6.3-bin.zip

# In order to correctly link we need to build tinyalsa ourself
RUN git clone https://github.com/tinyalsa/tinyalsa.git
WORKDIR /tinyalsa
# Prepare our header files
RUN make
RUN make install
RUN ldconfig
# Now we need to build the shared library (libtinyalsa.so
WORKDIR /tinyalsa/src
# In order to link it correctly we modify the Makefile a bit
# (We need a library with any version number after .so)
RUN sed -i 's\-shared -Wl,-soname,\-shared -soname=\g' Makefile
RUN sed -i 's\libtinyalsa.so.$(LIBVERSION_MAJOR)\libtinyalsa.so\g' Makefile
RUN sed -i 's\libtinyalsa.so.$(LIBVERSION)\libtinyalsa.so\g' Makefile
RUN sed -i 's\ln -sf\#ln -sf\g' Makefile
RUN make clean
# Now we compile it using the Android NDK
RUN make CC=$CC LD=$LD AR=$AR AS=$AS libtinyalsa.a libtinyalsa.so

# Lets arrange some libraries and header files
RUN cp -r /usr/local/include/tinyalsa/ /usr/include/
RUN cp -r /usr/include/tinyalsa/ $NDK_INCLUDE
RUN cp /tinyalsa/src/libtinyalsa.so $NDK_LIBS/libtinyalsa.so

# Finished preparing our compiler