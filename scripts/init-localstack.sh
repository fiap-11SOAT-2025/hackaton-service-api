#!/bin/bash
echo "⏳ Iniciando configuração do LocalStack..."

# Criar Bucket S3
awslocal s3 mb s3://fiap-videos || true
echo "✅ Bucket 'fiap-videos' criado."

# Criar Fila SQS
awslocal sqs create-queue --queue-name video-processing-queue || true
echo "✅ Fila 'video-processing-queue' criada."

# Configuração de SNS
# Criar o tópico SNS que o Worker espera
awslocal sns create-topic --name video-notifications

# Criar uma subscrição de e-mail para testes (opcional no LocalStack)
awslocal sns subscribe \
    --topic-arn arn:aws:sns:us-east-1:000000000000:video-notifications \
    --protocol email \
    --notification-endpoint no-reply@fiapx.com

# Criar Secret no Secrets Manager
awslocal secretsmanager create-secret \
    --name db-credentials \
    --secret-string '{"DB_HOST":"db","DB_USERNAME":"user","DB_PASSWORD":"password","DB_NAME":"fiapx_db","DB_SSL_MODE":"disable","DB_PORT":"5432"}'

echo "✅ LocalStack configurado com S3, SQS, SNS e Secrets Manager!"