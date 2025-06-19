# 🧪 Testinator

**Testinator** is an open-source framework that lets you describe and run backend test cases in **natural language**, powered by LLMs.

Write tests like:

> "Create a user → Check the database → Verify a Kafka event → Ensure the audit service was called"

...and Testinator will take care of the rest.

---

## 🔍 Features

- 💬 **Natural language test cases**  
  Describe tests using plain English or a lightweight DSL.

- 📄 **Contract-aware**  
  Parses OpenAPI, Protobuf, and Avro schemas to understand your API and message formats.

- 🧠 **LLM-powered reasoning**  
  Uses a language model to generate actionable test plans from your descriptions.

- 🧪 **Full execution support**  
  - Sends real HTTP/gRPC requests  
  - Verifies database state (PostgreSQL)  
  - Listens to Kafka topics and validates messages  
  - Inspects external service calls (via Wiremock)

- 📈 **Readable test reports**  
  Get clear logs and results for each test step.

---

## 💡 Example

```yaml
name: Create user and audit
steps:
  - POST /users with payload { "name": "Alice" }
  - Verify in DB: users table has entry with name "Alice"
  - Expect Kafka topic "user.created" to contain message with "id"
  - Check that /audit was called on external service
