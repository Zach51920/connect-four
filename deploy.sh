#!/bin/bash

# Navigate to the project directory
cd ~/connect-four

# Stop the existing process if it's running
pkill -f "connect-four"

# Ensure the .env file exists and update the COOKIE_SECRET
if [ -f .env ]; then
    # If .env exists, update the COOKIE_SECRET
    sed -i '/COOKIE_SECRET/d' .env
else
    # If .env doesn't exist, create it from example.env
    cp example.env .env
fi

# Add or update COOKIE_SECRET
echo "COOKIE_SECRET=$COOKIE_SECRET" >> .env

# Start the new version of the application
nohup ./bin/connect-four > app.log 2>&1 &

echo "Deployment completed. Application is running."

# Optionally, you can add Nginx configuration here if needed
# sudo cp nginx.conf /etc/nginx/conf.d/connect-four.conf
# sudo nginx -t && sudo systemctl reload nginx
