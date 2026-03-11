#!/bin/bash
# AegisGate Production Deployment Script
# Auto-generated deployment script

echo 'Starting AegisGate production deployment...'

# Deployment commands here
kubectl apply -f ./kubernetes/

echo 'Deployment complete.'
