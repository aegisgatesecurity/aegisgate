# AegisGate LangChain Integration

Comprehensive security integration for LangChain applications with real-time threat detection, compliance monitoring, and security filtering.

## Overview

The AegisGate LangChain integration provides enterprise-grade security for LLM applications built with LangChain. It acts as a security middleware layer that intercepts all LLM interactions to detect and block threats in real-time.

## Features

### Threat Detection
- **Prompt Injection Prevention** - Detect and block LLM01攻击 vectors
- **PII/PHI Detection** - Automatic identification and redaction of sensitive data
- **Toxicity Filtering** - Content safety monitoring for inappropriate content
- **Secret Detection** - Identify exposed API keys, tokens, and credentials

### Compliance Monitoring
- **MITRE ATLAS** - Malicious LLM Threats framework compliance
- **NIST AI RMF** - AI Risk Management Framework alignment
- **OWASP LLM** - OWASP Top 10 for LLM applications
- **SOC2, HIPAA, GDPR, ISO 27001** - Multi-framework support

### Security Filters
- `AegisGatePromptInjectionFilter`
- `AegisGatePIIFilter`
- `AegisGateToxicityFilter`
- `AegisGateSecretFilter`

## Quick Start

```python
from aegisgate.integrations.langchain import AegisGateCallbackHandler
from langchain_community.llms import OpenAI

# Initialize AegisGate security handler
aegisgate_handler = AegisGateCallbackHandler(
    api_key="your_api_key",
    endpoint="https://api.aegisgate.io",
    block_mode=True
)

# Use with LangChain
llm = OpenAI(
    callbacks=[aegisgate_handler],
    temperature=0.7
)

# All interactions are automatically protected
response = llm.predict("Hello, how are you?")
```

## Callback Handler

The `AegisGateCallbackHandler` class provides comprehensive security monitoring for all LangChain operations:

```python
from aegisgate.integrations.langchain import AegisGateCallbackHandler

handler = AegisGateCallbackHandler(
    api_key="your_api_key",
    endpoint="https://api.aegisgate.io",
    block_mode=True,
    compliance_frameworks=["SOC2", "HIPAA", "GDPR"]
)
```

### Methods

- `on_chain_start()` - Monitor chain execution
- `on_llm_start()` - Monitor LLM interactions
- `on_chain_end()` - Log chain completion
- `on_llm_end()` - Log LLM completion
- `on_chain_error()` - Handle chain errors
- `on_llm_error()` - Handle LLM errors

## Configuration

```python
from aegisgate.integrations.langchain import AegisGateConfig

config = AegisGateConfig(
    api_key="your_api_key",
    endpoint="https://api.aegisgate.io",
    block_mode=True,
    log_level="INFO",
   redact_pii=True,
    redact_phi=True
)
```

## Security Filters

### Prompt Injection Filter
Detects and blocks prompt injection attacks:

```python
from aegisgate.integrations.langchain import AegisGatePromptInjectionFilter

filter = AegisGatePromptInjectionFilter(threshold=0.8)
```

### PII Filter
Identifies and redacts personally identifiable information:

```python
from aegisgate.integrations.langchain import AegisGatePIIFilter

filter = AegisGatePIIFilter(
    entities=["PERSON", "EMAIL", "PHONE", "SSN"],
    redact=True
)
```

### Toxicity Filter
Monitors for inappropriate content:

```python
from aegisgate.integrations.langchain import AegisGateToxicityFilter

filter = AegisGateToxicityFilter(threshold=0.7)
```

### Secret Filter
Detects exposed credentials:

```python
from aegisgate.integrations.langchain import AegisGateSecretFilter

filter = AegisGateSecretFilter(patterns="auto")
```

## Compliance Frameworks

Support for multiple compliance frameworks:

| Framework | Status | Description |
|-----------|--------|-------------|
| MITRE ATLAS | ✅ | Malicious LLM Threats framework |
| NIST AI RMF | ✅ | AI Risk Management Framework |
| OWASP LLM | ✅ | OWASP Top 10 for LLMs |
| SOC2 | ✅ | SOC2 Type II compliance |
| HIPAA | ✅ | HIPAA data protection |
| GDPR | ✅ | GDPR privacy compliance |
| ISO 27001 | ✅ | Information security management |

## Integration Examples

### With Chat Models

```python
from langchain_community.chat_models import ChatOpenAI
from aegisgate.integrations.langchain import AegisGateChatModelWrapper

# Wrap chat model with AegisGate security
llm = ChatOpenAI(model="gpt-4")
wrapped_llm = AegisGateChatModelWrapper(llm, api_key="your_api_key")

response = wrapped_llm.predict("User question with PII: SSN 123-45-6789")
```

### With Chains

```python
from langchain.chains import LLMChain
from langchain.prompts import PromptTemplate
from aegisgate.integrations.langchain import AegisGateCallbackHandler

prompt = PromptTemplate.from_template("Answer: {question}")
llm = ChatOpenAI()

handler = AegisGateCallbackHandler(api_key="your_api_key")

chain = LLMChain(
    llm=llm,
    prompt=prompt,
    callbacks=[handler]
)

result = chain.run(question="What is the capital of France?")
```

## Testing

The integration includes comprehensive test coverage:

```bash
cd aegisgate/integrations/langchain
pytest tests/
```

## API Reference

### AegisGateCallbackHandler

```python
class AegisGateCallbackHandler:
    def __init__(
        self,
        api_key: str,
        endpoint: str,
        block_mode: bool = True,
        compliance_frameworks: List[str] = None,
        log_level: str = "INFO"
    ):
        """Initialize AegisGate security handler."""
    
    def on_chain_start(self, serialized: Dict[str, Any], inputs: Dict[str, Any], **kwargs):
        """Called when a chain starts."""
    
    def on_llm_start(self, serialized: Dict[str, Any], prompts: List[str], **kwargs):
        """Called when LLM starts generating."""
    
    def on_chain_end(self, outputs: Dict[str, Any], **kwargs):
        """Called when a chain ends."""
    
    def on_llm_end(self, response: LLMResult, **kwargs):
        """Called when LLM finishes generating."""
```

### AegisGateChatModelWrapper

```python
class AegisGateChatModelWrapper:
    def __init__(self, chat_model: BaseChatModel, config: AegisGateConfig = None):
        """Wrap a LangChain chat model with AegisGate security."""
    
    def predict(self, text: str) -> str:
        """Predict with security filtering."""
    
    def predict_messages(self, messages: List[BaseMessage]) -> BaseMessage:
        """Predict with messages and security filtering."""
```

## Error Handling

```python
from aegisgate.exceptions import AegisGateError, SecurityBlockError

try:
    response = wrapped_llm.predict("User input")
except SecurityBlockError as e:
    print(f"Request blocked: {e.reason}")
except AegisGateError as e:
    print(f"AegisGate error: {e}")
```

## Performance

- **Latency**: <5ms per request
- **Throughput**: 50,000+ requests/second
- **Memory**: ~128MB base footprint

## Support

For issues and questions:
- GitHub Issues: https://github.com/aegisgatesecurity/aegisgate/issues
- Documentation: https://aegisgate.io/docs
- Email: support@aegisgatesecurity.io

## License

MIT License - Copyright (c) 2025-2026 AegisGate Security
