# Dockerfile for the Svelte frontend (Vite)
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json* ./
# Ensure devDependencies (svelte-kit CLI) are installed during build
ENV NODE_ENV=development
RUN npm ci --legacy-peer-deps

FROM node:20-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

FROM nginx:stable-alpine AS runtime
# SvelteKit outputs client assets to .svelte-kit/output/client when using svelte-kit build
COPY --from=build /app/.svelte-kit/output/client /usr/share/nginx/html
COPY ../nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["/usr/sbin/nginx","-g","daemon off;"]
