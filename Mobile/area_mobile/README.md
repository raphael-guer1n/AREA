
The AREA mobile app is a thin Flutter client for browsing services, authenticating users, and launching OAuth flows on the go. It talks directly to the backend APIs exposed by the gateway.

## Highlights
- **Authentication:** Email/password login and signup, plus OAuth (Google/Apple/Facebook) via backend-issued URLs; tokens stored securely with `flutter_secure_storage`.
- **Service catalog:** Fetches available providers for the signed-in user, launches external browser OAuth, and listens for deep-link callbacks to refresh connection status.
- **Bottom-nav shell:** Dashboard, AREA placeholder, Services, and Profile tabs under a single `MainShell` scaffold.
- **Profile:** Displays basic user metadata and supports logout.

## Architecture (what’s in the code)
- **Framework:** Flutter with Material 3 styling (`lib/theme`).
- **State:** `Provider` (`AuthProvider`) drives authentication state, loading, and error feedback.
- **Networking:** `http` package + `flutter_dotenv` for the API base URL (`BASE_URL`).
- **Auth flow:** `AuthService` handles email/password endpoints and generic OAuth flows; listens to deep links via `app_links` with redirect URI `com.example.area_mobile:/oauth2redirect`.
- **Services flow:** `ServiceConnector` hits `/oauth2/providers/{userId}` and `/oauth2/authorize` to start provider connections; disconnect is stubbed for now.

## Screens (lib/screens)
- `auth/login_screen.dart` — Email/username + password login, Google OAuth button, validation, inline error banner.
- `auth/register_screen.dart` — Signup form with validation and immediate navigation into the main shell.
- `main_shell.dart` — Bottom navigation hosting:
  - `home/home_screen.dart` (dashboard placeholder)
  - `area/area_screen.dart` (automation placeholder)
  - `services/services_screen.dart` (provider list, search, connect/disconnect UI)
  - `profile/profile_screen.dart` (user info + logout)

## Project Structure
```
lib/
├── main.dart          # App entry, Provider setup, routing
├── providers/         # AuthProvider (session state + secure storage)
├── services/          # AuthService, ServiceConnector (HTTP + OAuth helpers)
├── screens/           # UI screens (auth, services, profile, shell)
├── models/            # ServiceModel (provider data)
└── theme/             # Colors, typography, theme
```

## Setup & Run
```bash
cd Mobile/area_mobile
echo "BASE_URL=http://10.0.2.2:8080/auth-service" > .env   # point to the gateway
flutter pub get
flutter run                      # choose your emulator/device
```
Notes:
- `BASE_URL` should be the gateway origin (e.g., `http://10.0.2.2:8080/auth-service` for the Android emulator).
- OAuth redirect is hardcoded to `com.example.area_mobile:/oauth2redirect`; keep Android/iOS configs aligned.

## Testing & Release
- Run tests: `flutter test`
- Build artifacts: `flutter build apk --release` or `flutter build ios --release`

## Known Gaps
- Service disconnect currently throws `UnimplementedError` in `ServiceConnector.disconnectService`; add a backend endpoint + hook it up to enable removal of linked accounts.
