# KubeCortex
A Vendor-Neutral Protocol for Constitutional Agentic Control Planes. Defining the Cognitive Agent Interface (CAI) for Kubernetes.
# KubeCortex
> **A Vendor-Neutral Protocol for Constitutional Agentic Control Planes.**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/kubecortex)](https://goreportcard.com/report/github.com/YOUR_USERNAME/kubecortex)
[![CNCF Landscape](https://img.shields.io/badge/CNCF%20Landscape-Sandbox-blue)](https://landscape.cncf.io/) *<-- (Only add "Sandbox" once accepted, use "Candidate" for now)*

---

**KubeCortex** defines the **Cognitive Agent Interface (CAI)**, a standard protocol that allows any AI Agent (AutoGen, LangGraph, CrewAI) to safely negotiate with Kubernetes clusters via a "Constitutional Proxy."

It solves the **"Black Box"** problem in Agentic Infrastructure by wrapping probabilistic AI reasoning in deterministic safety policies (OPA/Kyverno).

## 🏗 Architecture
* **Protocol:** CAI (Cognitive Agent Interface)
* **Safety:** The Constitutional Proxy (Sidecar)
* **State:** Pluggable Durable Backend (Default: Temporal)

flowchart TB
    %% =====================
    %% Entry Point
    %% =====================
    User[User / CI-CD System<br/>PlanRequest CRD]:::core
    Operator[Cortex Operator<br/>(Semantic Router)]:::core

    User --> Operator

    %% =====================
    %% Layer 1: Cognitive Planner
    %% =====================
    subgraph L1["Layer 1: Cognitive Planner (Probabilistic Brain)"]
        Router[Semantic Router]:::core

        subgraph Agents["Pluggable AI Planners (CAI Protocol)"]
            AutoGen[AutoGen<br/>Multi-Agent Debate]:::ai
            LangGraph[LangGraph<br/>Stateful Reasoning]:::ai
            CrewAI[CrewAI<br/>Role-Based Agents]:::ai
        end

        Draft[Draft Plan<br/>(YAML Manifest)]:::ai
    end

    Operator --> Router
    Router --> AutoGen
    Router --> LangGraph
    Router --> CrewAI

    AutoGen --> Draft
    LangGraph --> Draft
    CrewAI --> Draft

    %% =====================
    %% Layer 2: Constitutional Proxy
    %% =====================
    subgraph L2["Layer 2: Constitutional Proxy (Deterministic Firewall)"]
        Proxy[Constitutional Proxy<br/>(API Isolation)]:::safety
        Policy[Policy Engine<br/>(OPA / Kyverno)]:::safety
        Reject[Policy Violation<br/>Structured Error]:::reject
        Sign[Verified & Signed Plan]:::exec
    end

    Draft --> Proxy
    Proxy --> Policy
    Policy -- Reject --> Reject
    Reject -. Self-Correction Loop .-> Router
    Policy -- Approve --> Sign

    %% =====================
    %% Layer 3: Durable Execution
    %% =====================
    subgraph L3["Layer 3: Durable Execution (The Hands)"]
        Executor[Pluggable Executor]:::exec
        Temporal[Temporal Worker<br/>Durable Workflow]:::exec
        Argo[Argo Workflows<br/>DAG Executor]:::exec
    end

    Sign --> Executor
    Executor --> Temporal
    Executor --> Argo

    %% =====================
    %% Layer 4: Infrastructure
    %% =====================
    subgraph L4["Layer 4: Infrastructure (The Target)"]
        K8s[Kubernetes API]:::infra
        Pods[Pods / Nodes]:::infra
        Crossplane[Crossplane Resources<br/>(Cloud Infra)]:::infra
    end

    Temporal --> K8s
    Argo --> K8s
    K8s --> Pods
    K8s --> Crossplane

    %% =====================
    %% State Persistence
    %% =====================
    subgraph Memory["Durable State & Audit Trail"]
        Store[Postgres / Redis<br/>Immutable History]:::memory
    end

    Operator --> Store
    Proxy --> Store
    Executor --> Store

    %% =====================
    %% Styles
    %% =====================
    classDef core fill:#1e3a8a,color:#fff,stroke:#60a5fa
    classDef ai fill:#6d28d9,color:#fff,stroke:#c4b5fd
    classDef safety fill:#7f1d1d,color:#fff,stroke:#f87171
    classDef reject fill:#991b1b,color:#fff,stroke:#ef4444
    classDef exec fill:#065f46,color:#fff,stroke:#34d399
    classDef infra fill:#374151,color:#fff,stroke:#9ca3af
    classDef memory fill:#3f3f46,color:#fff,stroke:#a1a1aa


## 🚀 Quick Start
KubeCortex is designed to run on any Kubernetes cluster (Kind, Minikube, EKS, GKE).

```bash
# 1. Install the KubeCortex Operator
helm repo add kubecortex [https://charts.kubecortex.io](https://charts.kubecortex.io)
helm install cortex-operator kubecortex/operator

# 2. Submit your first Intent
kubectl apply -f examples/intent-high-availability.yaml
