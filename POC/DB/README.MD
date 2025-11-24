# Database POC Comparison: PostgreSQL, MongoDB, Redis

## Project Context

Three backend Proofs of Concept (POCs) were developed for AREA to compare different database technologies and determine which one is the most suitable for the project’s architecture and long-term needs.

Each POC implements the same minimal functional flow:

* Create or initialize a database/collection/key-space.
* Insert a user document/row (email, first_name, last_name).
* Retrieve the list of users or the stored user value.

The goals of this comparative study were to evaluate:

* Data modeling flexibility
* Relational integrity and structure
* Performance for typical AREA queries
* Ease of setup, deployment, and collaboration
* Scalability and reliability
* Developer experience and learning curve

---

## Comparative Table (Summary)

| Criteria                      | PostgreSQL (SQL)                                             | MongoDB (NoSQL Document Store)                                      | Redis (In-Memory Key-Value Store)                                  |
|------------------------------|---------------------------------------------------------------|-----------------------------------------------------------------------|---------------------------------------------------------------------|
| Data model                   | Strongly relational, normalized schemas                      | Flexible JSON documents, schema-less                                 | Simple key-value pairs (hashes, lists, sets)                        |
| Relations                    | Native FK, JOIN support                                      | Emulated via references, manual handling                             | Not supported                                                       |
| Query capabilities           | Very powerful SQL, indexing, transactions                    | Rich document queries, aggregation pipeline                          | Extremely limited                                                   |
| JSON support                 | Excellent (JSONB with indexing)                              | Native & flexible                                                     | Possible but not practical                                           |
| Performance                 | High for structured queries                                  | High for document access                                              | Extremely fast (in-memory)                                          |
| Persistence                  | Full ACID, durable storage                                   | Durable but less strict                                               | Optional, primarily RAM-based                                        |
| Setup & collaboration        | Requires schema & migrations, but predictable                | No schema required, auto-created collections                         | Minimal setup, no schema                                            |
| Use cases                    | Complex models, integrity constraints                        | Flexible or rapidly evolving data structures                         | Cache, ephemeral data, job queues                                   |
| Suitability for AREA         | Excellent: matches relational nature of AREA data            | Partial: good for settings/configs, not for core AREA relationships   | Poor as primary DB, good as auxiliary cache                          |

---

## Advantages & Drawbacks

### PostgreSQL (Selected Final Technology)

**Strengths**

* Clear, structured, and consistent relational modeling.
* Ideal for entities with strong dependencies (users, services, actions, reactions, areas).
* Full ACID guarantees and transactional safety.
* Powerful support for semi-structured data via JSONB.
* Predictable, mature, and well-understood tooling.
* Fits perfectly with the backend modular monolithic architecture.

**Limitations**

* Requires schema definitions and migrations.
* Slightly more rigid during rapid prototyping compared to NoSQL solutions.

---

### MongoDB

**Strengths**

* Very flexible schema: no migrations required.
* Native JSON storage ideal for dynamic or rapidly changing objects.
* Easy to get started, especially for simple POCs.
* Good performance for retrieving and updating documents.

**Limitations**

* No built-in relations: foreign keys, joins, and integrity rely on the backend code.
* More risk of data inconsistency without strong validation.
* Less suited for structured, relational-heavy data models.

*In the context of AREA, MongoDB would work for flexible configurations (e.g., dynamic service settings), but not for the core data model.*

---

### Redis

**Strengths**

* Extremely fast read/write operations (RAM-based).
* Perfect for caching, sessions, rate limiting, or temporary state.
* Minimal setup and operational overhead.

**Limitations**

* Not designed for permanent, structured, or relational data.
* Limited querying capabilities.
* Data persistence is optional and less reliable than SQL/NoSQL databases.

*In AREA, Redis could be useful as a complementary component (caching, task queues), but cannot serve as the main operational datastore.*

---

## Conclusion (Selected Database)

PostgreSQL was selected as the primary database for AREA because it provides the best combination of:

* Strong relational modeling  
* Data consistency and integrity  
* Ability to store both structured and semi-structured data (JSONB)  
* Mature tooling, indexing, and query capabilities  
* Predictability and maintainability for a long-term project  
* Seamless integration with the backend architecture  

While MongoDB and Redis bring valuable strengths in their respective domains, they do not meet the structural and consistency requirements needed for AREA’s interconnected entities (users → services → actions → reactions → areas).

PostgreSQL offers the reliability and flexibility required for the final application while maintaining a straightforward development workflow for the entire team.
