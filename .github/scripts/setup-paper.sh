#!/usr/bin/env bash
#
# Downloads PaperMC, enables the management server, starts it, and waits until
# it is ready. Every CI job that needs a live server calls this so the three of
# them cannot drift apart.
#
# With TLS=true it also generates a self-signed certificate for the server and
# exports the public half as a PEM, so the test client can trust it without
# turning verification off. The certificate names both localhost and 127.0.0.1,
# which lets either be used as TEST_HOST.
#
# Required: PAPERMC_VERSION, TEST_SERVER_PORT, TEST_SERVER_SECRET
# Optional: TLS (default false), KEYSTORE_PASSWORD, SERVER_DIR (default server)

set -euo pipefail

: "${PAPERMC_VERSION:?PAPERMC_VERSION is required}"
: "${TEST_SERVER_PORT:?TEST_SERVER_PORT is required}"
: "${TEST_SERVER_SECRET:?TEST_SERVER_SECRET is required}"

TLS="${TLS:-false}"
KEYSTORE_PASSWORD="${KEYSTORE_PASSWORD:-changeit}"
SERVER_DIR="${SERVER_DIR:-server}"

mkdir -p "$SERVER_DIR"
cd "$SERVER_DIR"

java -version

echo "Fetching latest PaperMC build for ${PAPERMC_VERSION}..."
PAPER_URL=$(curl -s -A "mcrpc-github-ci/1.0.0 (https://github.com/fi-xz/mcrpc)" \
  "https://fill.papermc.io/v3/projects/paper/versions/${PAPERMC_VERSION}/builds/latest" |
  jq -r '.downloads."server:default".url')

if [ -z "$PAPER_URL" ] || [ "$PAPER_URL" = "null" ]; then
  echo "Failed to fetch PaperMC download URL"
  exit 1
fi

echo "Downloading PaperMC from: $PAPER_URL"
curl -o paper.jar "$PAPER_URL"

echo "eula=true" >eula.txt

TLS_ENABLED=false
TLS_KEYSTORE=""
TLS_PASSWORD=""

if [ "$TLS" = "true" ]; then
  echo "Generating a self-signed certificate for the management server..."
  keytool -genkeypair \
    -alias mcrpc \
    -keyalg RSA -keysize 2048 \
    -storetype PKCS12 \
    -keystore keystore.p12 \
    -storepass "$KEYSTORE_PASSWORD" \
    -dname "CN=localhost, O=mcrpc CI" \
    -ext "SAN=dns:localhost,ip:127.0.0.1" \
    -validity 1

  # The client trusts this PEM as a root, so the run exercises real certificate
  # verification rather than skipping it.
  keytool -exportcert \
    -alias mcrpc \
    -keystore keystore.p12 \
    -storepass "$KEYSTORE_PASSWORD" \
    -rfc -file ca.crt

  TLS_ENABLED=true
  TLS_KEYSTORE=keystore.p12
  TLS_PASSWORD="$KEYSTORE_PASSWORD"
fi

cat >server.properties <<EOF
management-server-allowed-origins=
management-server-enabled=true
management-server-host=localhost
management-server-port=${TEST_SERVER_PORT}
management-server-secret=${TEST_SERVER_SECRET}
management-server-tls-enabled=${TLS_ENABLED}
management-server-tls-keystore=${TLS_KEYSTORE}
management-server-tls-keystore-password=${TLS_PASSWORD}
EOF

nohup java -Xms512M -Xmx1G -jar paper.jar nogui >server.log 2>&1 &
echo $! >server.pid

for i in $(seq 1 60); do
  if grep -q "Done" server.log 2>/dev/null; then
    echo "Server started successfully"
    exit 0
  fi
  if [ "$i" -eq 60 ]; then
    echo "Server failed to start within 60 seconds"
    cat server.log
    exit 1
  fi
  sleep 1
done
