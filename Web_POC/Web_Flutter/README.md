# Flutter Web Frontend POC

## POC Overview
This Flutter Web prototype replicates a mini-AREA flow: homepage, login form, and profile page displaying entered credentials.

## Included Features
- Homepage with navigation to the login page.
- Login page with basic validation (email/password) and error message.
- Profile page displaying the email and password submitted via the navigation.
- Light/dark theme managed by Flutter (Material 3) and simple responsive layout.

## Test Credentials
- Email: `test@test.com`
- Password: `0000`

## Installation Instructions (Linux)
1. Install Flutter 3.5+ (SDK) and system dependencies (Chrome or a compatible web browser).
2. Check the environment:
```bash
flutter doctor
```
3. Retrieve project dependencies:
```bash
flutter pub get
```

## Commands to launch the project
- Start locally (Chrome):
```bash
flutter run -d chrome
```
- Browserless version (web server):
```bash
flutter run -d web-server --web-hostname localhost --web-port 8080
```
- Static production build:
```bash
flutter build web
```

## Strengths of Flutter Web
* Single Dart code that can be used for web, mobile, and desktop.
* Hot reload and Flutter tooling enable rapid iteration.
* Material 3 widgets and centralized theming (colors, fonts, styles).
* Declarative navigation via named routes, easy to implement.

## Limitations observed (specific to the POC)
* Intentionally simulated client-side authentication: no API or persistent storage; state disappears upon refresh.
* Credentials are stored and displayed only for demonstration purposes; not suitable for production.
* No global state management, navigation middleware, or automated tests yet, which could be added in a full version.
* Performance and SEO need to be validated on more demanding scenarios than this prototype.

## Conclusion
The POC validates navigation, data entry, and information display in Flutter Web. Further development: integrate an authentication API, secure transport/storage (token/session), add tests (widget/integration), and measure performance in production (using a Flutter Web build and static deployment).