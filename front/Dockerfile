FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* pnpm-lock.yaml* ./

RUN npm install -g pnpm && pnpm install

COPY . .

RUN pnpm build

FROM nginx:alpine AS runner

WORKDIR /usr/share/nginx/html

# Clean default nginx static files
RUN rm -rf ./*

COPY --from=builder /app/dist .
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]

