#!/bin/sh
# Script para esperar a que PostgreSQL esté listo

set -e

echo "Esperando a que PostgreSQL esté disponible..."

# Intentar conectarse a PostgreSQL usando variables de entorno
until PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" -c "SELECT 1;" > /dev/null 2>&1; do
  echo "PostgreSQL no está disponible aún - esperando..."
  sleep 1
done

echo "PostgreSQL está listo!"