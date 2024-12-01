FROM node:lts-slim

# Install necessary dependencies for Electron
RUN apt-get update && apt-get install -y \
    libgtk-3-0 \
    libx11-xcb1 \
    libxcb-dri3-0 \
    libxss1 \
    libnss3 \
    libasound2 \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application
COPY . .

# Expose the port your app uses for HTTP traffic
EXPOSE 3000

# Command to run your Electron app
CMD ["npm", "run", "run-build"]
