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

```mermaid
flowchart TB

    User["User / CI-CD System\nPlanRequest CRD"]
    Operator["Cortex Operator\nSemantic Router"]

    User --> Operator

    subgraph L1[Layer 1: Cognitive Planner]
        Router["Semantic Router"]

        subgraph Agents[Pluggable AI Planners]
            AutoGen["AutoGen\nMulti-Agent Debate"]
            LangGraph["LangGraph\nStateful Reasoning"]
            CrewAI["CrewAI\nRole-Based Agents"]
        end

        Draft["Draft Plan\nYAML Manifest"]
    end

    Operator --> Router
    Router --> AutoGen
    Router --> LangGraph
    Router --> CrewAI

    AutoGen --> Draft
    LangGraph --> Draft
    CrewAI --> Draft

    subgraph L2[Layer 2: Constitutional Proxy]
        Proxy["Constitutional Proxy"]
        Policy["Policy Engine\nOPA / Kyverno"]
        Reject["Policy Violation"]
        Sign["Verified & Signed Plan"]
    end

    Draft --> Proxy
    Proxy --> Policy
    Policy -->|Reject| Reject
    Reject -.-> Router
    Policy -->|Approve| Sign

    subgraph L3[Layer 3: Durable Execution]
        Executor["Pluggable Executor"]
        Temporal["Temporal Worker"]
        Argo["Argo Workflows"]
    end

    Sign --> Executor
    Executor --> Temporal
    Executor --> Argo

    subgraph L4[Layer 4: Infrastructure]
        K8s["Kubernetes API"]
        Pods["Pods / Nodes"]
        Crossplane["Crossplane Resources"]
    end

    Temporal --> K8s
    Argo --> K8s
    K8s --> Pods
    K8s --> Crossplane

    subgraph Memory[State & Audit Trail]
        Store["Postgres / Redis"]
    end

    Operator --> Store
    Proxy --> Store
    Executor --> Store
```

## 🚀 Quick Start
KubeCortex is designed to run on any Kubernetes cluster (Kind, Minikube, EKS, GKE).

```bash
# 1. Install the KubeCortex Operator
helm repo add kubecortex [https://charts.kubecortex.io](https://charts.kubecortex.io)
helm install cortex-operator kubecortex/operator

# 2. Submit your first Intent
kubectl apply -f examples/intent-high-availability.yaml
