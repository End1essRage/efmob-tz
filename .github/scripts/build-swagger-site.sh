#!/bin/bash
# .github/scripts/build-swagger-site.sh
# Сборка Swagger документации для GitHub Pages

set -e

echo "Building Swagger site..."

# Создаем директории
mkdir -p public/docs

# Копируем Swagger UI
cp -r node_modules/swagger-ui-dist/* public/

# Копируем спецификации
IFS=',' read -ra SERVICE_LIST <<< "$SERVICES"
for service in "${SERVICE_LIST[@]}"; do
    SWAGGER_FILE="docs/${service}-swagger.json"
    if [ ! -f "$SWAGGER_FILE" ]; then
        echo "❌ $SWAGGER_FILE not found"
        exit 1
    fi
    cp "$SWAGGER_FILE" "public/docs/${service}-swagger.json"
done

# Создаем основной index.html
cat > public/index.html << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>API Documentation - Subs Service</title>
    <link rel="stylesheet" type="text/css" href="./swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *,
        *:before,
        *:after {
            box-sizing: inherit;
        }
        body {
            margin: 0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="./swagger-ui-bundle.js"></script>
    <script src="./swagger-ui-standalone-preset.js"></script>
    <script>
    window.onload = function() {
        window.ui = SwaggerUIBundle({
            url: "./docs/subs-swagger.json",
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            plugins: [
                SwaggerUIBundle.plugins.DownloadUrl
            ],
            layout: "StandaloneLayout"
        });
    };
    </script>
</body>
</html>
EOF

echo "Swagger site built successfully."