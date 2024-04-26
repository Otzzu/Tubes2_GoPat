# Use a specific version of node on Alpine for better predictability
FROM node:21-alpine3.19

# Create a new group and user 'app' for running our application
RUN addgroup -S app && adduser -S -G app app

# Set the working directory
WORKDIR /app

# Copy package files and ensure permissions are correct before installing dependencies
COPY package*.json ./
RUN chown -R app:app /app

# Switch to user 'app' before installing node modules
USER app
RUN npm install

# Copy the rest of the application's code
COPY --chown=app:app . .

# Build the application as user 'app'
RUN npm run build

# The app binds to port 3000, make sure the container exposes it
EXPOSE 3000

# Command to run the app
CMD ["npm", "start"]
