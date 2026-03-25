# Dockerfile for the Svelte frontend (Vite)
FROM oven/bun:1-alpine AS deps
WORKDIR /app
COPY package.json bun.lock* ./
# Ensure devDependencies (svelte-kit CLI) are installed during build
ENV NODE_ENV=development
RUN bun install --frozen-lockfile

FROM oven/bun:1-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN bun run build

FROM nginx:stable-alpine AS runtime
# SvelteKit with adapter-static outputs to 'build/' by default
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["/usr/sbin/nginx","-g","daemon off;"]
