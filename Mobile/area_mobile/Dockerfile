FROM ghcr.io/cirruslabs/flutter:latest

WORKDIR /app

# Base URL injected into the generated .env (used by scripts/build-apk.sh)
ARG BASE_URL=http://gateway:8080
ENV BASE_URL=${BASE_URL}

# Where the APK will be copied inside the container (shared via volume)
ENV APK_OUTPUT_PATH=/apk/client.apk
ENV EXIT_AFTER_BUILD=true

# Prime dependencies early for better layer caching
COPY pubspec.* ./
RUN flutter pub get

# Copy the full project
COPY . .

# Ensure build script is executable
RUN chmod +x scripts/build-apk.sh

# Regenerate local.properties with container SDK paths to avoid host-specific values
RUN FLUTTER_ROOT="$(dirname "$(dirname "$(which flutter)")")" && \
    ANDROID_SDK_ROOT="${ANDROID_HOME:-/opt/android-sdk}" && \
    echo "flutter.sdk=${FLUTTER_ROOT}" > android/local.properties && \
    echo "sdk.dir=${ANDROID_SDK_ROOT}" >> android/local.properties && \
    echo "flutter.buildMode=release" >> android/local.properties && \
    echo "flutter.versionName=1.0.0" >> android/local.properties && \
    echo "flutter.versionCode=1" >> android/local.properties && \
    echo "flutter.compileSdkVersion=34" >> android/local.properties && \
    echo "flutter.minSdkVersion=19" >> android/local.properties && \
    echo "flutter.targetSdkVersion=34" >> android/local.properties

# Accept Android licenses non-interactively (best effort)
RUN yes | flutter doctor --android-licenses || true

# Clean and prepare a fresh build
RUN flutter clean

# Build the APK and copy it to /apk (volume-mounted in compose)
CMD ["sh", "-c", "./scripts/build-apk.sh"]
