<div align="center">
  <img height="300" src="https://github.com/kokaq/.github/blob/main/kokaq-protocol.png" alt="cute quokka as kokaq logo"/>
</div>

`kokaq-protocol` contains the canonical protocol definitions for the kokaq distributed priority queue system. It serves as the single source of truth for how components in the `kokaq` ecosystem communicate — ensuring consistency, versioning, and interoperability across clients, servers, and storage backends.

## 🧠 Why Protocol Definitions?
Distributed systems require a clear and consistent language to operate reliably across environments, runtimes, and network boundaries. kokaq-protocol ensures:

- ✅ Stable Contracts — Breaking changes are avoided through strict versioning
- 🔄 Multi-language Support — Generate code for Go, Rust, Python, C#, etc.
- 📡 Standardized Communication — Between clients, servers, and storage engines
- 🔍 Evolvable Design — Protocols can evolve while remaining backward-compatible

## 🔗 Related Projects
- `kokaq-core` — Priority queue logic and scheduling
- `kokaq-server` — Runtime system using this protocol
- `kokaq-client` — SDKs that consume this protocol
