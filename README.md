<div align="center">
    <h1>AREA Project</h1>
    <h3>Automation Platform Inspired by IFTTT and Zapier</h3>
</div>

<br>

The AREA project aims to develop a complete platform that allows users to automate actions between different online services.

The concept is simple: **if an Action occurs in one service, then a REAction is automatically executed in another service.**

<br>

## Project Objective

The goal of this project is to build a software suite composed of three main components:

### 1. **Application Server**

The server centralizes all business logic.  
It manages users, external services, actions, reactions, and the automatic execution of automations (called *AREAs*).  
It exposes a REST API used by the web and mobile clients.

---

### 2. **Web Application**

The web application provides a graphical user interface to interact with the platform.

It allows users to:
- create an account and authenticate,
- connect to external services,
- configure Actions and ReActions,
- create and manage Areas,
- view and organize existing automations.

No business logic is executed on the web client: it communicates exclusively with the server.

---

### 3. **Mobile Application**

The mobile application provides the same features as the web interface, optimized for smartphones.

It allows users to:
- access their account,
- manage subscribed services,
- create and manage AREAs.

Just like the web client, no business logic is executed on the mobile app.

<br>

## What is an AREA?

An **AREA** is an automation composed of two elements:

- **Action** — An event detected in a service (e.g., "Email received", "New file created", "New social post").
- **REAction** — A task automatically executed when the Action occurs (e.g., "Send a message", "Store a file", "Publish a post").

This system allows users to connect multiple services and automate digital tasks seamlessly.

<br>

## The Role of the Hook

The server uses a system of *Hooks* to detect when user-defined Actions occur.  
When a Hook identifies that an Action has been triggered, it automatically launches the associated REAction using the server API.

<br>

## Documentation & Ressources

Cette section regroupe toutes les documentations, schémas et outils utilisés pour concevoir et structurer le projet AREA.

### Documents Techniques
- Documentation officielle du projet (PDF) — [`G-DEV-500_AREA.pdf`][Area-Subject]

### Design & UX
- Maquettes Figma — [Template][figma-template]

### Diagrammes & Architecture
- Schémas Mermaid (Architecture, Séquence, Classes) — [Diagrammes Mermaid][mermaid-diagrams]

[figma-template]: https://www.figma.com/make/hAuXyYuX12okDvawHiHnPX/Template-for-IFTTT-App?node-id=0-4&t=sbehe7KY4mvlfYbj-1
[mermaid-diagrams]: https://www.mermaidchart.com/app/projects/d32e2d39-2142-4e20-84c5-20cbb151cc1d/diagrams/e0382344-f52b-4dc9-8ced-bdec15998fb4/share/invite/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkb2N1bWVudElEIjoiZTAzODIzNDQtZjUyYi00ZGM5LThjZWQtYmRlYzE1OTk4ZmI0IiwiYWNjZXNzIjoiVmlldyIsImlhdCI6MTc2NDA5OTEzNH0.eFtTAe5t2KfZcl2UWPqwu_YgcjjlUX2CZpOwfZl31dA
[Area-Subject]: https://intra.epitech.eu/module/2025/G-DEV-500/PAR-5-2/acti-692707/project/file/G-DEV-500_AREA.pdf

<br>

## Educational Objectives

The AREA project aims to teach:
- how to integrate multiple external services and APIs,
- how to design and structure a complete software architecture,
- the separation of interface and business logic,
- the internal functioning of modern automation platforms.

<br>

## Methodology

- Development following an **Agile methodology**, using 2-week sprints  
- **TDD** (Test-Driven Development) for core functionalities  
- **CI/CD** pipelines through GitHub Actions for automated testing and deployment  

<br>

## Authors

- Alexandre Guillaud – Developer  
- Alexis Constantinopoulos – Developer  
- Raphaël Guerin – Developer  
- Clément-Alexis Fournier – Developer  

<br>

## License – MIT License

MIT License

Copyright (c) 2025 AREA Project Team

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

<br>

<div align="center"> <sub>{Epitech} — 2025</sub> </div>
