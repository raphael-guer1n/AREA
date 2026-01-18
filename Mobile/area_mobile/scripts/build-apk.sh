#!/bin/sh
set -e

APK_OUTPUT_PATH="${APK_OUTPUT_PATH:-/apk/client.apk}"
APK_OUTPUT_DIR="$(dirname "$APK_OUTPUT_PATH")"
BASE_URL="${BASE_URL:-http://gateway:8080}"

echo "-> Generating .env for BASE_URL=$BASE_URL"
printf "BASE_URL=%s\n" "$BASE_URL" > .env

echo "-> Ensuring dependencies are installed"
flutter pub get

echo "-> Building release APK"
flutter build apk --release

APK_SOURCE="build/app/outputs/flutter-apk/app-release.apk"
if [ ! -f "$APK_SOURCE" ]; then
  echo "APK not found at $APK_SOURCE" >&2
  exit 1
fi

echo "-> Copying APK to shared volume: $APK_OUTPUT_PATH"
mkdir -p "$APK_OUTPUT_DIR"
cp "$APK_SOURCE" "$APK_OUTPUT_PATH"
chmod 755 "$APK_OUTPUT_DIR"
chmod 644 "$APK_OUTPUT_PATH"

echo "APK ready at $APK_OUTPUT_PATH"
if [ "${EXIT_AFTER_BUILD:-false}" = "true" ]; then
  exit 0
fi

exec tail -f /dev/null
