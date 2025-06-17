#!/bin/sh
# Function to safely read a secret file with fallback paths
read_secret() {
  local secret_name=$1
  local default_value=${2:-""}

  # Try Docker Swarm path first
  if [ -f "/run/secrets/${secret_name}" ]; then
    cat "/run/secrets/${secret_name}"
    echo "Found secret at /run/secrets/${secret_name}" >&2
  # Try Docker Compose v3.x path
  elif [ -f "/${secret_name}" ]; then
    cat "/${secret_name}"
    echo "Found secret at /${secret_name}" >&2
  # Try Bitnami container path
  elif [ -f "/opt/bitnami/scripts/secrets/${secret_name}" ]; then
    cat "/opt/bitnami/scripts/secrets/${secret_name}"
    echo "Found secret at /opt/bitnami/scripts/secrets/${secret_name}" >&2
  # Try Docker Compose alternative path
  elif [ -f "/secrets/${secret_name}" ]; then
    cat "/secrets/${secret_name}"
    echo "Found secret at /secrets/${secret_name}" >&2
  # Try local secrets directory (for development)
  elif [ -f "/app/secrets/${secret_name}.txt" ]; then
    cat "/app/secrets/${secret_name}.txt"
    echo "Found secret at /app/secrets/${secret_name}.txt" >&2
  # Try local secrets directory (alternative path)
  elif [ -f "./secrets/${secret_name}.txt" ]; then
    cat "./secrets/${secret_name}.txt"
    echo "Found secret at ./secrets/${secret_name}.txt" >&2
  # Check if the secret is available as an environment variable
  elif [ -n "$(eval echo \$$(echo "$secret_name" | tr '[:lower:]' '[:upper:]'))" ]; then
    eval echo \$$(echo "$secret_name" | tr '[:lower:]' '[:upper:]')
    echo "Found secret in environment variable ${secret_name}" >&2
  # Return default value if file not found
  else
    # For Redis, provide a default password if ALLOW_EMPTY_PASSWORD is set
    if [ "${secret_name}" = "redis_password" ] && [ "${ALLOW_EMPTY_PASSWORD}" = "yes" ]; then
      echo ""
      echo "Using empty password for Redis due to ALLOW_EMPTY_PASSWORD=yes" >&2
    else
      # Use default values for MongoDB if not found
      if [ "${secret_name}" = "mongo_root_username" ]; then
        echo "root"
        echo "Using default value 'root' for mongo_root_username" >&2
      elif [ "${secret_name}" = "mongo_root_password" ]; then
        echo "password"
        echo "Using default value 'password' for mongo_root_password" >&2
      else
        echo "${default_value}"
        echo "Warning: Secret ${secret_name} not found at expected paths" >&2
        # Print all environment variables for debugging
        echo "Available environment variables:" >&2
        env | grep -i secret || true
        # List contents of common secret directories
        echo "Contents of /run/secrets/ directory:" >&2
        ls -la /run/secrets/ 2>/dev/null || echo "Directory not found or empty" >&2
        echo "Contents of / directory (looking for secret files):" >&2
        ls -la / | grep -i secret 2>/dev/null || echo "No secret files found" >&2
        echo "Contents of /app/secrets/ directory:" >&2
        ls -la /app/secrets/ 2>/dev/null || echo "Directory not found or empty" >&2
        echo "Contents of ./secrets/ directory:" >&2
        ls -la ./secrets/ 2>/dev/null || echo "Directory not found or empty" >&2
      fi
    fi
  fi
}

# Set environment variables using the read_secret function
export GF_SECURITY_ADMIN_PASSWORD=$(read_secret "gf_security_admin_password")
export GF_SECURITY_ADMIN_USER=$(read_secret "gf_security_admin_user")
export MONGODB_ROOT_PASSWORD=$(read_secret "mongo_root_password")
export MONGODB_ROOT_USERNAME=$(read_secret "mongo_root_username")
export POSTGRESQL_PASSWORD=$(read_secret "postgresql_password")
export POSTGRESQL_USERNAME=$(read_secret "postgresql_username")
export REDIS_PASSWORD=$(read_secret "redis_password")

# Debug: Print environment variables
#echo "Environment variables set:"
#echo "MONGODB_ROOT_USERNAME: $MONGODB_ROOT_USERNAME"
#echo "MONGODB_ROOT_PASSWORD: $MONGODB_ROOT_PASSWORD"
#echo "POSTGRESQL_USERNAME: $POSTGRESQL_USERNAME"
#echo "POSTGRESQL_PASSWORD: $POSTGRESQL_PASSWORD"

exec "$@"
