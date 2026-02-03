#!/bin/bash
echo "⏳ Iniciando configuração do LocalStack..."

# Criar Bucket S3
awslocal s3 mb s3://fiap-videos || true
echo "✅ Bucket 'fiap-videos' criado."

# Criar Fila SQS
awslocal sqs create-queue --queue-name video-processing-queue || true
echo "✅ Fila 'video-processing-queue' criada."

# 3. Configurar E-mail (SES)
# Verifica o remetente oficial do sistema (obrigatório na AWS)
awslocal ses verify-email-identity --email-address no-reply@fiapx.com
echo "✅ Emails verificados"
awslocal ses list-identities
echo "✅ LocalStack configurado: SES pronto para envio!"