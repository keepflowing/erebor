FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci --omit=dev

COPY . .

RUN npm run build

# second stage
FROM nginx:alpine

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /app /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
