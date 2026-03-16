# AegisGate Architecture Diagrams

This document contains architecture diagrams for AegisGate using Mermaid.

---

## 1. High-Level Architecture

```mermaid
flowchart TB
    subgraph Clients["Client Applications"]
        direction LR
        WebApp["Web App"]
        MobileApp["Mobile App"]
        Backend["Backend Service"]
        AIAgent["AI Agent"]
    end

    subgraph AegisGate["AegisGate Proxy"]
        direction TB
        TLS["TLS Termination"]

        subgraph SecurityLayer["Security Layer"]
            Auth["Auth & Rate Limit"]
            RBAC["RBAC"]
            Secrets["Secret Rotation"]
        end

        subgraph ComplianceLayer["Compliance Layer"]
            OWASP["OWASP Checker"]
            SOC2["SOC2"]
            HIPAA["HIPAA"]
            PCI["PCI-DSS"]
        end

        subgraph MLLayer["ML Detection Layer"]
            Anomaly["Anomaly Detection"]
            Threat["Threat Intel"]
        end

        subgraph Observability["Observability"]
            Metrics["Metrics"]
            Logging["Logging"]
            Audit["Audit Logger"]
        end

        Proxy["Request Proxy"]
    end

    subgraph Providers["AI Providers"]
        OpenAI["OpenAI"]
        Anthropic["Anthropic"]
        Cohere["Cohere"]
        Azure["Azure OpenAI"]
    end

    subgraph Data["Data Layer"]
        Postgres["PostgreSQL"]
        Redis["Redis"]
        File["File Storage"]
    end

    Clients --> TLS
    TLS --> Auth
    Auth --> ComplianceLayer
    ComplianceLayer --> MLLayer
    MLLayer --> Proxy
    Proxy --> Providers
    Proxy --> Observability
    Observability --> Data
```

---

## 2. Request Flow

```mermaid
sequenceDiagram
    participant Client as Client App
    participant AegisGate as AegisGate Proxy
    participant ML as ML Detector
    participant Compliance as Compliance
    participant Provider as AI Provider

    Client->>AegisGate: 1. Request

    Note over AegisGate: 2. TLS Termination
    Note over AegisGate: 3. Auth Check (API Key/JWT)

    AegisGate->>AegisGate: 4. Rate Limiter Check
    alt Rate Limited
        AegisGate-->>Client: 429 Too Many Requests
    end

    AegisGate->>Compliance: 5. Compliance Check
    alt Violation Found
        Compliance->>AegisGate: Log & Alert
        AegisGate-->>Client: 403 Forbidden
    end

    AegisGate->>ML: 6. ML Analysis
    alt Threat Detected
        ML->>AegisGate: Block & Alert
        AegisGate-->>Client: 403 Forbidden
    end

    AegisGate->>Provider: 7. Forward Request
    Provider-->>AegisGate: AI Response
    AegisGate-->>Client: 8. Response

    Note over AegisGate: 9. Log & Metrics
```

---

## 3. Component Architecture

```mermaid
flowchart TB
    subgraph API["API Layer"]
        HTTP[":8080 HTTP"]
        HTTPS[":8443 HTTPS"]
        GRPC[":9001 gRPC"]
    end

    subgraph Proxy["Proxy Engine"]
        Route["Route Matcher"]
        Transform["Request Transform"]
        Adapter["Provider Adapter"]
        LB["Load Balancer"]
    end

    subgraph Security["Security"]
        TLS["TLS/mTLS"]
        Auth["Authentication"]
        Rate["Rate Limiter"]
        JWT["JWT Validator"]
    end

    subgraph Compliance["Compliance Engine"]
        Registry["Framework Registry"]
        Scanner["Pattern Scanner"]
        Reporter["Finding Reporter"]
    end

    subgraph ML["ML Detection"]
        Traffic["Traffic Analyzer"]
        Anomaly["Anomaly Detector"]
        Threat["Threat Intel"]
    end

    subgraph Storage["Storage"]
        DB["PostgreSQL"]
        Cache["Redis"]
        Files["File Storage"]
    end

    HTTP --> Proxy
    HTTPS --> Proxy
    GRPC --> Proxy

    Proxy --> Security
    Security --> Compliance
    Compliance --> ML
    ML --> Adapter

    Adapter --> DB
    Adapter --> Cache
    Adapter --> Files
```

---

## 4. Tier Architecture

```mermaid
flowchart LR
    subgraph Community["COMMUNITY (Free)"]
        C1["Basic Proxy"]
        C2["Rate Limiter"]
        C3["OWASP View"]
        C4["Basic Metrics"]
    end

    subgraph Developer["DEVELOPER ($29/mo)"]
        D1["All Community"]
        D2["PostgreSQL"]
        D3["OAuth/SSO"]
        D4["Cost Alerts"]
        D5["Grafana"]
    end

    subgraph Professional["PROFESSIONAL ($99/mo)"]
        P1["All Developer"]
        P2["HIPAA/PCI/SOC2"]
        P3["Multi-Tenant"]
        P4["ML Threat Detection"]
        P5["SIEM Integration"]
    end

    subgraph Enterprise["ENTERPRISE (Custom)"]
        E1["All Professional"]
        E2["ISO 42001/FedRAMP"]
        E3["HSM/FIPS"]
        E4["On-Premise"]
        E5["Air-Gapped"]
    end

    Community --> Developer
    Developer --> Professional
    Professional --> Enterprise
```

---

## 5. Kubernetes Deployment

```mermaid
flowchart TB
    subgraph K8s["Kubernetes Cluster"]
        LB["Load Balancer"]

        subgraph AegisGateDeploy["AegisGate Deployment"]
            Pod1["Pod 1"]
            Pod2["Pod 2"]
            PodN["Pod N"]
        end

        subgraph Services["Services"]
            Svc["AegisGate Service"]
        end

        subgraph Data["Data Layer"]
            PG["PostgreSQL\n(RDS)"]
            Redis["Redis\n(Cluster)"]
            S3["S3 Bucket"]
        end

        subgraph Monitoring["Monitoring"]
            Graf["Grafana"]
            Prometheus["Prometheus"]
        end
    end

    LB --> Svc
    Svc --> Pod1
    Svc --> Pod2
    Svc --> PodN

    Pod1 --> PG
    Pod2 --> Redis
    PodN --> S3

    Pod1 --> Graf
    Pod2 --> Prometheus
```

---

## 6. Database Schema Overview

```mermaid
erDiagram
    Tenant ||--o{ User : "has"
    Tenant ||--o{ APIKey : "issues"
    Tenant ||--o{ Route : "configures"
    Tenant ||--o{ AuditLog : "generates"

    User ||--o{ Session : "creates"

    Route ||--o{ ComplianceCheck : "triggers"
    Route ||--o{ RequestLog : "logs"

    APIKey ||--o{ RequestLog : "authenticates"

    Tenant {
        string id PK
        string name
        string tier
        timestamp created_at
    }

    User {
        string id PK
        string tenant_id FK
        string email
        string role
        timestamp created_at
    }

    APIKey {
        string id PK
        string tenant_id FK
        string key_hash
        string name
        timestamp expires_at
        timestamp created_at
    }

    Route {
        string id PK
        string tenant_id FK
        string path
        string target
        string provider
    }

    RequestLog {
        string id PK
        string route_id FK
        string api_key_id FK
        string method
        string path
        int status_code
        int latency_ms
        timestamp created_at
    }

    AuditLog {
        string id PK
        string tenant_id FK
        string action
        string actor
        json details
        timestamp created_at
    }
```

---

## 7. Rate Limiting Flow

```mermaid
flowchart TB
    Request["Request"] --> KeyGen["Generate Key\n(API Key or IP)"]

    KeyGen --> Lookup["Lookup in\nToken Bucket"]

    subgraph Check["Rate Limit Check"]
        Tokens["Tokens Available?"]
    end

    Lookup --> Tokens

    alt Yes
        Tokens --> Decrement["Decrement Token"]
        Decrement --> Forward["Forward Request"]
        Forward --> Success["200 OK"]
    end

    alt No
        Tokens --> Block["Block Request"]
        Block --> RateLimited["429 Rate Limited"]
        RateLimited --> Retry["Retry-After: 60s"]
    end

    subgraph Cleanup["Background Cleanup"]
        Timer["Timer\n(every minute)"]
        Timer --> Remove["Remove Stale Entries"]
    end
```

---

## 8. Compliance Checking Flow

```mermaid
flowchart TB
    Request["Request"] --> Parse["Parse Request"]

    Parse --> Match["Match Against\nPatterns"]

    subgraph Evaluate["Evaluation"]
        OWASP["OWASP Patterns"]
        PII["PII Detection"]
        Sensitive["Sensitive Data"]
    end

    Match --> Evaluate

    Evaluate --> Violation["Violation Found?"]

    alt Yes
        Violation --> Log["Log Finding"]
        Log --> Alert["Generate Alert"]
        Alert --> Action["Action: Block/Log/Warn"]
    end

    alt No
        Violation --> Forward["Forward Request"]
    end

    subgraph Report["Reporting"]
        Summary["Generate Summary"]
        Export["Export: JSON/PDF"]
    end

    Log --> Report
    Forward --> Report
```

---

## 9. License Activation Flow

```mermaid
sequenceDiagram
    participant User as User
    participant AegisGate as AegisGate
    participant License as License Server

    User->>AegisGate: Enter License Key
    AegisGate->>License: Validate License

    alt Valid License
        License-->>AegisGate: License Valid
        AegisGate->>AegisGate: Enable Tier Features
        AegisGate-->>User: Success! Tier Activated
    end

    alt Invalid License
        License-->>AegisGate: Invalid/Expired
        AegisGate-->>User: Error: Invalid License
    end

    alt Community (No License)
        AegisGate-->>AegisGate: Use Community Features
        AegisGate-->>User: Running Community Tier
    end
```

---

## 10. Security Layers

```mermaid
flowchart TB
    subgraph L1["Layer 1: Network"]
        TLS["TLS 1.3"]
        Firewall["Firewall"]
    end

    subgraph L2["Layer 2: Transport"]
        MTLS["mTLS"]
        Certs["Certificates"]
    end

    subgraph L3["Layer 3: Application"]
        Auth["AuthN/AuthZ"]
        RBAC["RBAC"]
        JWT["JWT Validation"]
    end

    subgraph L4["Layer 4: Data"]
        Encrypt["Encryption at Rest"]
        Mask["Data Masking"]
        Audit["Audit Logging"]
    end

    L1 --> L2
    L2 --> L3
    L3 --> L4
```

---

*Diagrams rendered with Mermaid. For more information, see [architecture.md](architecture.md).*
