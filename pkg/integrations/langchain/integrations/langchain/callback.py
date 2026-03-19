# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Callback handlers for AegisGate LangChain integration.

This module provides LangChain callback handlers that integrate with 
AegisGate for monitoring and analyzing LLM application behavior.
"""

from __future__ import annotations

import logging
import time
from typing import Any, Dict, List, Optional, cast
from datetime import datetime

from langchain_core.callbacks import BaseCallbackHandler
from langchain_core.messages import BaseMessage, HumanMessage, AIMessage
from langchain_core.outputs import LLMResult, ChatResult
from langchain_core.outputs import ChatGeneration, Generation

from aegisgate.integrations.langchain.config import AegisGateConfig
from aegisgate.integrations.langchain.exceptions import (
    LangChainIntegrationError,
    ComplianceViolationError,
    RateLimitExceededError,
)
from aegisgate import AegisGateClient, ClientConfig

logger = logging.getLogger(__name__)


class AegisGateCallbackHandler(BaseCallbackHandler):
    """
    LangChain callback handler for AegisGate integration.
    
    This handler monitors LLM calls and events, sending relevant data
    to AegisGate for threat detection, compliance monitoring, and analysis.
    
    Attributes:
        config: AegisGate configuration
        client: AegisGate client instance
        start_times: Track start times for latency calculations
    """
    
    def __init__(
        self,
        config: Optional[AegisGateConfig] = None,
        client: Optional[AegisGateClient] = None,
    ):
        """
        Initialize the callback handler.
        
        Args:
            config: AegisGate configuration (optional, uses env vars if not provided)
            client: AegisGate client (optional, creates one if not provided)
        """
        super().__init__()
        
        self.config = config or AegisGateConfig.from_env()
        self.client = client
        self._client_initialized = False
        
        # Track start times for latency
        self.start_times: Dict[str, float] = {}
        
        logger.debug("AegisGateCallbackHandler initialized")
    
    @property
    def aegisgate_client(self) -> AegisGateClient:
        """Get or create the AegisGate client."""
        if not self._client_initialized and self.config.api_key:
            self.client = AegisGateClient(
                base_url=self.config.base_url,
                api_key=self.config.api_key,
                timeout=self.config.timeout,
            )
            self._client_initialized = True
        return self.client  # type: ignore
    
    def _send_event(self, event_type: str, data: Dict[str, Any]) -> None:
        """Send an event to AegisGate SIEM."""
        if not self.config.logging.enabled:
            return
        
        if not self.aegisgate_client:
            logger.debug(f"SIEM event {event_type} would be sent to AegisGate")
            return
        
        try:
            # Create SIEM event
            siem_event = {
                "event_type": event_type,
                "timestamp": datetime.utcnow().isoformat(),
                "source": "langchain",
                "data": data,
                "severity": data.get("severity", "info"),
            }
            
            # Send to SIEM
            self.aegisgate_client.siem.send_event(siem_event)
            logger.debug(f"Sent SIEM event: {event_type}")
            
        except Exception as e:
            logger.warning(f"Failed to send SIEM event: {e}")
    
    def on_chain_start(
        self,
        serialized: Dict[str, Any],
        inputs: Dict[str, Any],
        *,
        run_id: str,
        parent_run_id: Optional[str] = None,
        tags: Optional[List[str]] = None,
        metadata: Optional[Dict[str, Any]] = None,
        **kwargs: Any,
    ) -> None:
        """Called when a chain starts."""
        self.start_times[run_id] = time.time()
        
        if not self.config.logging.enabled:
            return
        
        logger.debug(f"Chain start: {run_id}")
        
        # Log chain start to AegisGate
        self._send_event(
            "chain_start",
            {
                "chain_name": serialized.get("name", "unknown"),
                "inputs": inputs,
                "run_id": run_id,
                "parent_run_id": parent_run_id,
            },
        )
    
    def on_chain_end(
        self,
        outputs: Dict[str, Any],
        *,
        run_id: str,
        parent_run_id: Optional[str] = None,
        tags: Optional[List[str]] = None,
        **kwargs: Any,
    ) -> None:
        """Called when a chain ends."""
        start_time = self.start_times.pop(run_id, None)
        
        if not self.config.logging.enabled:
            return
        
        # Calculate duration
        duration = time.time() - start_time if start_time else None
        
        logger.debug(f"Chain end: {run_id}, duration: {duration}s")
        
        # Log chain end to AegisGate
        self._send_event(
            "chain_end",
            {
                "outputs": outputs,
                "run_id": run_id,
                "duration_seconds": duration,
            },
        )
    
    def on_llm_start(
        self,
        serialized: Dict[str, Any],
        prompts: List[str],
        *,
        run_id: str,
        parent_run_id: Optional[str] = None,
        tags: Optional[List[str]] = None,
        metadata: Optional[Dict[str, Any]] = None,
        **kwargs: Any,
    ) -> None:
        """Called when LLM starts generating."""
        self.start_times[run_id] = time.time()
        
        if not self.config.logging.enabled:
            return
        
        # Log prompt to AegisGate for threat detection
        if self.aegisgate_client and self.config.threat_detection.enabled:
            try:
                # Combine all prompts for threat detection
                combined_prompt = "\n\n".join(prompts)
                
                # Check for threats in prompts
                result = self.aegisgate_client.proxy.check_threats(
                    content=combined_prompt,
                    content_type="input",
                )
                
                if result.get("threat_detected"):
                    threat_type = result.get("threat_type", "unknown")
                    threat_score = result.get("threat_score", 0)
                    
                    # Log the threat
                    self._send_event(
                        "threat_detected",
                        {
                            "threat_type": threat_type,
                            "threat_score": threat_score,
                            "run_id": run_id,
                        },
                    )
                    
                    # Apply configuration action
                    if self.config.threat_action == "block":
                        raise LangChainIntegrationError(
                            message=f"Threat detected during LLM call: {threat_type}",
                            threat_type=threat_type,
                            threat_score=threat_score,
                        )
                        
            except Exception as e:
                logger.error(f"Error during LLM prompt threat detection: {e}")
        
        logger.debug(f"LLM start: {run_id}, prompts: {len(prompts)}")
        
        # Log LLM start
        self._send_event(
            "llm_start",
            {
                "model_name": serialized.get("name", "unknown"),
                "prompts_count": len(prompts),
                "run_id": run_id,
            },
        )
    
    def on_llm_end(
        self,
        response: LLMResult,
        *,
        run_id: str,
        parent_run_id: Optional[str] = None,
        tags: Optional[List[str]] = None,
        **kwargs: Any,
    ) -> None:
        """Called when LLM finishes generating."""
        start_time = self.start_times.pop(run_id, None)
        duration = time.time() - start_time if start_time else None
        
        if not self.config.logging.enabled:
            return
        
        # Log output for PII detection
        if self.aegisgate_client and self.config.pii_redaction.enabled:
            try:
                # Extract generated text
                if response.generations:
                    generated_text = response.generations[0][0].text
                    
                    # Check for PII in output
                    result = self.aegisgate_client.proxy.check_threats(
                        content=generated_text,
                        content_type="output",
                    )
                    
                    if result.get("threat_detected"):
                        self._send_event(
                            "pii_detected",
                            {
                                "threat_type": "pii",
                                "threat_score": result.get("threat_score", 0),
                                "run_id": run_id,
                            },
                        )
                        
            except Exception as e:
                logger.warning(f"Error during LLM output checking: {e}")
        
        logger.debug(f"LLM end: {run_id}, duration: {duration}s")
        
        # Log LLM end
        self._send_event(
            "llm_end",
            {
                "generations": len(response.generations),
                "token_usage": response.llm_output.get("token_usage") if response.llm_output else None,
                "run_id": run_id,
                "duration_seconds": duration,
            },
        )
    
    def on_tool_start(
        self,
        serialized: Dict[str, Any],
        input_str: str,
        *,
        run_id: str,
        parent_run_id: Optional[str] = None,
        tags: Optional[List[str]] = None,
        metadata: Optional[Dict[str, Any]] = None,
        **kwargs: Any,
    ) -> None:
        """Called when a tool starts."""
        self.start_times[run_id] = time.time()
        
        if not self.config.logging.enabled:
            return
        
        logger.debug(f"Tool start: {run_id}")
        
        # Log tool start
        self._send_event(
            "tool_start",
            {
                "tool_name": serialized.get("name", "unknown"),
                "input": input_str,
                "run_id": run_id,
            },
        )
    
    def on_tool_end(
        self,
        output: str,
        *,
        run_id: str,
        parent_run_id: Optional[str] = None,
        tags: Optional[List[str]] = None,
        **kwargs: Any,
    ) -> None:
        """Called when a tool ends."""
        start_time = self.start_times.pop(run_id, None)
        duration = time.time() - start_time if start_time else None
        
        if not self.config.logging.enabled:
            return
        
        logger.debug(f"Tool end: {run_id}, duration: {duration}s")
        
        # Log tool end
        self._send_event(
            "tool_end",
            {
                "output": output,
                "run_id": run_id,
                "duration_seconds": duration,
            },
        )
    
    def on_error(
        self,
        error: Exception,
        *,
        run_id: str,
        parent_run_id: Optional[str] = None,
        tags: Optional[List[str]] = None,
        **kwargs: Any,
    ) -> None:
        """Called when an error occurs."""
        start_time = self.start_times.pop(run_id, None)
        duration = time.time() - start_time if start_time else None
        
        logger.error(f"Error in run {run_id}: {type(error).__name__}: {error}")
        
        # Log error to AegisGate
        self._send_event(
            "error",
            {
                "error_type": type(error).__name__,
                "error_message": str(error),
                "run_id": run_id,
                "duration_seconds": duration,
            },
        )
    
    def on_retriever_start(
        self,
        serialized: Dict[str, Any],
        query: str,
        *,
        run_id: str,
        parent_run_id: Optional[str] = None,
        tags: Optional[List[str]] = None,
        metadata: Optional[Dict[str, Any]] = None,
        **kwargs: Any,
    ) -> None:
        """Called when retriever starts."""
        self.start_times[run_id] = time.time()
        
        if not self.config.logging.enabled:
            return
        
        logger.debug(f"Retriever start: {run_id}")
        
        # Log retriever start
        self._send_event(
            "retriever_start",
            {
                "query": query,
                "run_id": run_id,
            },
        )
    
    def on_retriever_end(
        self,
        documents: List[Any],
        *,
        run_id: str,
        parent_run_id: Optional[str] = None,
        tags: Optional[List[str]] = None,
        **kwargs: Any,
    ) -> None:
        """Called when retriever ends."""
        start_time = self.start_times.pop(run_id, None)
        duration = time.time() - start_time if start_time else None
        
        if not self.config.logging.enabled:
            return
        
        logger.debug(f"Retriever end: {run_id}, docs: {len(documents)}")
        
        # Log retriever end
        self._send_event(
            "retriever_end",
            {
                "documents_count": len(documents),
                "run_id": run_id,
                "duration_seconds": duration,
            },
        )
