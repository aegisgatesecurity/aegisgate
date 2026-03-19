# AegisGate

[![Go version](https://img.shields.io/badge/Go-1.25%2B-orange.svg)](https://golang.org/)
[![Version](https://img.shields.io/github/v/release/aegisgatesecurity/aegisgate)](https://github.com/aegisgatesecurity/aegisgate/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

AegisGate is an Enterprise AI API Security Platform that secures AI API traffic between client applications and AI service providers (OpenAI, Anthropic, Azure, AWS Bedrock, Cohere).

## Features

- **Zero code changes required** - Deploy as drop-in gateway
- **Real-time threat blocking** - Prompt injection, data leakage, adversarial attacks
- **Out-of-the-box compliance** - SOC2, HIPAA, PCI-DSS, GDPR, ISO 27001/42001, NIST AI RMF
- **ML-powered detection** - Advanced threat intelligence
- **Less than 5ms latency overhead** - Minimal performance impact

## Quick Links

- [Go Documentation](https://docs.aegisgatesecurity.io)
- [Python SDK Documentation](temp_aegisgate/sdk/python/README.md)
- [LangChain Integration](temp_aegisgate/sdk/python/aegisgate/integrations/langchain/)

## Installation (Go)

```bash
go mod download github.com/aegisgatesecurity/aegisgate
```

## Installation (Python)

```bash
pip install aegisgate[langchain-openai]
```

## Quick Start (Python with LangChain)

```python
from aegisgate.integrations.langchain import AegisGateConfig, AegisGateChatModelWrapper
from langchain_openai import ChatOpenAI

# Create your LLM
llm = ChatOpenAI(api_key="your-key")

# Wrap with AegisGate security
config = AegisGateConfig(
    base_url="http://localhost:8080",
    api_key="AG-your-key",
    threat_detection=True,
    block_on_critical=True,
)

secure_llm = AegisGateChatModelWrapper(llm=llm, config=config)

# Use normally - all calls are now monitored and secured
response = secure_llm.invoke([{"role": "user", "content": "Hello!"}])
```

## Development

### Build

```bash
go build -o aegisgate cmd/aegisgate/main.go
```

### Test

```bash
go test ./...
# or for Python SDK
cd temp_aegisgate/sdk/python && pytest tests/
```

## Support

- Bug Reports: https://github.com/aegisgatesecurity/aegisgate/issues
- Security Issues: security@aegisgatesecurity.io
- Documentation: https://docs.aegisgatesecurity.io

## License

MIT License - See [LICENSE](LICENSE) file for details
