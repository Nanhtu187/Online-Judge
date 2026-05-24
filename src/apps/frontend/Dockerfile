FROM node:20-alpine AS builder

WORKDIR /app

# Enable Corepack and install dependencies
RUN corepack enable
COPY src/apps/frontend ./
RUN yarn install
RUN yarn build

# Serve with Nginx
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
