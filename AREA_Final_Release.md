# 🚀 Final Release — AREA

## Overview

**AREA** is an automation platform inspired by IFTTT and Zapier, designed to connect services together through configurable **Actions** and **Reactions**.  
Its goal is to let users automate everyday workflows (professional or personal) without writing any code.

This final release represents the completion of the project from a **technical**, **functional**, and **product** perspective, with a full deployment on a **public domain**.

- 🌐 Frontend: https://app.area-connect.cloud/  
- ⚙️ Backend (API): https://api.area-connect.cloud/

---

## Product Vision

AREA is based on a simple concept:

> **When an Action happens on a service, one or more Reactions are automatically triggered on other services.**

### Concrete examples
- New Outlook email → Create a Google Calendar event  
- New YouTube video → Send a Discord message  
- Extreme weather alert → Send an email  
- New Dropbox file → Automatically copy it to OneDrive  

Users build these workflows using **AREAs**, without having to deal with the underlying technical complexity.

---

## Global Architecture

AREA relies on a **microservices architecture** designed to be scalable, reliable, and maintainable.

### Backend (Go)
- **API Gateway**: single entry point, routing, authentication, security policies  
- **AuthService**: user accounts, JWT, OAuth2  
- **ServiceService**: service catalog, actions, reactions, OAuth configurations  
- **AreaService**: AREA creation, storage, and orchestration  
- **WebhookService**: incoming webhook handling  
- **PollingService**: polling-based actions (RSS, weather, external APIs)  
- **CronService**: time-based triggers  
- **MailService**: email delivery  

Each service is isolated and communicates via REST / gRPC.  
Database layer is powered by **PostgreSQL**.

---

## Clients

### Web Client
- Built with **Next.js**
- Modern, responsive, and accessible UI
- Full management of services and AREAs

### Mobile Client
- Built with **Flutter** (Android)
- Consistent experience with the web client
- No business logic on the client side

---

## Services & Automations

AREA provides a rich ecosystem of connected services, combining **webhooks**, **polling**, and **time-based triggers**.

### Actions (examples)
- GitHub, Dropbox, Outlook Mail & Calendar, OneDrive  
- YouTube, RSS, NewsAPI, NASA  
- OpenWeatherMap  
- Timer (daily, weekly, monthly, delay)  
- Notion, IDFM Traffic, Coinbase  

### Reactions (examples)
- Send Discord messages  
- Create and manage calendar events  
- Send emails  
- File management (cloud storage)  
- YouTube playlist management  

➡️ More than **20 services**, **39 actions**, and **27 reactions** are available in this release.

---

## Orchestration & Performance

- Smart hook-based triggering system  
- Polling only for services actually used by the user  
- Event deduplication  
- Support for multiple reactions triggered by a single action  

This ensures good performance and scalability, even with many active AREAs.

---

## Security & Reliability

- Secure authentication (JWT, OAuth2)  
- Secure token storage  
- Clear separation of responsibilities between services  
- Docker & Docker Compose for reproducible environments  
- Backend CI/CD with automated tests and builds  

---

## Accessibility & User Experience

- Lighthouse accessibility score above 90  
- Clean and consistent UI across web and mobile  
- Built-in colorblind mode  
- User experience inspired by IFTTT and Zapier standards  

---

## Visuals

*(Screenshots will be added here: dashboard, AREA creation flow, action/reaction selection, mobile app, etc.)*

---

## Deployment

- Publicly accessible frontend and backend  
- Production-ready infrastructure  
- Architecture designed to easily add new services and features  

---

## Conclusion

This final release of AREA delivers:

- A solid and well-justified microservices architecture  
- A complete and functional automation platform  
- A polished user experience  
- A strong foundation for future evolution  

AREA is not just an academic project — it is a real automation platform, designed as a scalable and extensible product.
