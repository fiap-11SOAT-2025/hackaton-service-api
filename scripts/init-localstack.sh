#!/bin/bash
echo "⏳ Iniciando configuração do LocalStack..."

# Criar Bucket S3
awslocal s3 mb s3://fiap-videos
echo "✅ Bucket 'fiap-videos' criado."

# Criar Fila SQS
awslocal sqs create-queue --queue-name video-processing-queue
echo "✅ Fila 'video-processing-queue' criada."