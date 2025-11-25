# Blazor Frontend POC

## Quick Overview
Blazor WebAssembly (NET 8) POC to validate simple navigation between Home, Login, and Profile with lightweight client-side state management.

## POC Objective
- Demonstrate a simulated client-side authentication flow.
- Test navigation and route protection based on in-memory state.
- Evaluate initial usability before connecting to a real backend.

## Implemented Features
- **Home**: Welcome message and button to the login page.
- **Login**: Form validated via DataAnnotations, error display, conditional login.
- **Profile**: Automatic redirection if not authenticated, display of credentials stored in memory.

## Test Credentials
- Email: `test@test.com`
- Password: `0000`

## Installation and Runtime Guide
1) Prerequisites: .NET SDK 8.0 installed.
2) Restore and run in development mode:
```bash
dotnet restore
dotnet run
```
The application starts on http://localhost:5110 (profile `http` in `launchSettings.json`).
3) Hot-reload option:
```bash
dotnet watch run
```

## Strengths of Blazor
* WebAssembly execution with the ability to reuse C# code between front and back.
* Reusable components and built-in two-way data binding.
* Native validation support (DataAnnotations) and dependency injection included.
* Full .NET development tooling available (hot-reload, debugging, VS/VS Code support).

## Limitations observed (specific to the POC)
* Authentication intentionally simplified for the demonstration: no API calls or persistence, so session resets on refresh.
* Credentials stored in memory and displayed in the UI, not intended for real production usage.
* No network error handling or role/permission management implemented at this stage.
* No automated tests or CI setup yet, which could be added in a production context.


## Conclusion
The POC shows navigation, form validation, and local state management under Blazor WebAssembly. Next steps: connect an authentication API, secure storage (token/session), add tests and harden access management.