
The **AREA Mobile App** is the smartphone interface of the AREA platform — an automation system inspired by **IFTTT** and **Zapier**.  
It allows users to connect their favorite online services and automate actions between them, anytime and anywhere.

## 🧩 Overview

The mobile application is developed using **Flutter**, providing a smooth cross-platform experience on both **Android** and **iOS**.

It communicates directly with the **AREA Application Server** through a REST API to manage:
- user authentication,
- services connections,
- creation of Actions and ReActions,
- management of AREAs (automations).

All business logic and data handling are managed by the backend server — the mobile app acts as a **thin client** responsible for user interaction and UX.

## 🚀 Features

- **User Authentication**
  - Account creation and login via the AREA server  
  - Secure token-based authentication  

- **Service Management**
  - Connect and disconnect third-party services  
  - View connected services linked to your account  

- **AREA Management**
  - Create new *AREAs* (Action → Reaction workflows)  
  - Configure conditions and linked services  
  - View, edit, or delete existing AREAs  

- **Responsive UI/UX**
  - Modern Flutter UI optimized for phones and tablets  
  - Dark/Light theme support  

## 🛠️ Tech Stack

| Component | Technology |
|------------|-------------|
| Framework | Flutter (Dart) |
| State Management | Provider / Riverpod / BLoC *(depending on setup)* |
| Networking | `http` / `dio` |
| Authentication | OAuth 2.0 / Token-based auth |
| API | REST API (AREA Server) |
| Environment Handling | `flutter_dotenv` |
| CI/CD | GitHub Actions (build & test) |

## 📦 Project Structure

```text
area_mobile/
├── lib/
│   ├── main.dart                # Application entry point
│   ├── screens/                 # UI screens and navigation
│   ├── widgets/                 # Reusable UI components
│   ├── models/                  # Data models for Area, Action, Reaction, etc.
│   ├── providers/               # State providers / Controllers
│   ├── services/                # API calls & local storage logic
│   └── utils/                   # Constants, themes, helpers
├── assets/                      # Images, icons, translations
├── test/                        # Unit and widget tests
└── pubspec.yaml                 # Dependencies and project config
```

## ⚙️ Getting Started

### Prerequisites
- **Flutter SDK** (>=3.0.0)
- **Dart SDK**
- Access to the AREA **backend API** (running locally or remotely)
- Android Studio / VS Code configured with Flutter

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/AREA-Project/area_mobile.git
   cd area_mobile
   ```

2. **Install dependencies**
   ```bash
   flutter pub get
   ```

3. **Create an environment file**
   ```
   cp .env.example .env
   ```
   Configure the API base URL and other keys if needed:
   ```
   API_BASE_URL=https://api.area-project.com
   ```

4. **Run the app**
   ```bash
   flutter run
   ```

## 🧪 Testing

Run unit and widget tests:

```bash
flutter test
```

Or use the integrated GitHub Actions CI pipeline for continuous testing.

## 📱 Building for Production

To build release versions:

```bash
flutter build apk --release
flutter build ios --release
```

Build artifacts can then be deployed to the Google Play Store or Apple App Store.

## 👥 Contributors

- Alexandre Guillaud – Developer  
- Alexis Constantinopoulos – Developer  
- Raphaël Guerin – Developer  
- Clément-Alexis Fournier – Developer  

## 📄 License – MIT License

This project is licensed under the **MIT License**.  
See the [LICENSE](../LICENSE) file for details.
