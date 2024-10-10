# Use official Caddy image from Docker Hub
FROM caddy:alpine

# Copy Caddyfile into the container
COPY Caddyfile /etc/caddy/Caddyfile

# Expose port 5000 to match the Caddyfile
EXPOSE 5000

# Run Caddy
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
