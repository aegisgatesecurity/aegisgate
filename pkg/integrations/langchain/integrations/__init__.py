# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
LangChain integration for AegisGate.

This package provides seamless integration between LangChain and AegisGate,
allowing you to secure your LLM applications with automatic threat detection,
compliance monitoring, and policy enforcement.

Example:
    Basic usage with OpenAI:
    >>> from langchain_openai import OpenAI
    >>> from aegisgate.integrations.langchain import AegisGateLLMWrapper
    >>> 
    >>> # Create your LLM
    >>> llm = OpenAI(api_key="...")
    >>> 
    >>> # Wrap with AegisGate for security
    >>> secure_llm = AegisGateLLMWrapper(llm)
    >>> 
    >>> # Use normally - all calls are now monitored
    >>> response = secure_llm.invoke("Hello, world!")

    With full configuration:
    >>> from aegisgate.integrations.langchain import AegisGateChatModelWrapper
    >>> from langchain_openai import ChatOpenAI
    >>> 
    >>> config = AegisGateConfig(
    ...     base_url="http://localhost:8080",
    ...     api_key="AG-your-key",
    ...     threat_detection=True,
    ...     block_on_critical=True,
    ... )
    >>> 
    >>> chat = ChatOpenAI(model="gpt-4")
    >>> secure_chat = AegisGateChatModelWrapper(chat, config=config)
    >>> response = secure_chat.invoke([{"role": "user", "content": "..."}])

    For callbacks:
    >>> from aegisgate.integrations.langchain import AegisGateCallbackHandler
    >>> 
    >>> handler = AegisGateCallbackHandler(
    ...     base_url="http://localhost:8080",
    ...     api_key="AG-your-key"
    ... )
    >>> 
    >>> # Use with any LangChain chain
    >>> chain.invoke("...", config={"callbacks": [handler]})
"""

from aegisgate.integrations.langchain.config import AegisGateConfig
from aegisgate.integrations.langchain.wrapper import (
    AegisGateLLMWrapper,
    AegisGateChatModelWrapper,
)
from aegisgate.integrations.langchain.callback import AegisGateCallbackHandler
from aegisgate.integrations.langchain.filters import (
    PromptInjectionFilter,
    PIIFilter,
    ToxicityFilter,
)

__all__ = [
    "AegisGateConfig",
    "AegisGateLLMWrapper",
    "AegisGateChatModelWrapper",
    "AegisGateCallbackHandler",
    "PromptInjectionFilter",
    "PIIFilter",
    "ToxicityFilter",
]
