# Image Flutter Android (Docker Hub)
FROM cirruslabs/flutter:3.24.0-android

WORKDIR /app

# Adresse par défaut vers le gateway depuis un émulateur Android (10.0.2.2 = host)
ARG BASE_URL=http://10.0.2.2:8080/auth-service
ENV BASE_URL=${BASE_URL}

COPY . .

# Génère le fichier .env consommé par flutter_dotenv
RUN printf "BASE_URL=%s\n" "$BASE_URL" > .env

RUN flutter pub get
RUN flutter build apk --release

# Conteneur qui reste vivant pour pouvoir docker cp l'APK
CMD ["sleep", "infinity"]
