FROM oven/bun:1-alpine AS deps
WORKDIR /app
COPY code/frontend/package.json code/frontend/bun.lock* ./
RUN apk add --no-cache gcompat
ENV NODE_ENV=development
RUN bun install --frozen-lockfile

FROM oven/bun:1-alpine AS build
WORKDIR /app
RUN apk add --no-cache gcompat
COPY --from=deps /app/node_modules ./node_modules
COPY code/frontend/ .
# Copy the changelog into the static folder so SvelteKit sees it during build
COPY CHANGELOG.md ./static/CHANGELOG.md
ARG VITE_API_BASE
ARG VITE_WS_BASE
ENV VITE_API_BASE=$VITE_API_BASE
ENV VITE_WS_BASE=$VITE_WS_BASE
RUN bun run build

FROM nginx:stable-alpine AS runtime
COPY --from=build /app/build /usr/share/nginx/html
# Fix nginx config path for root context
COPY code/frontend/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["/usr/sbin/nginx", "-g", "daemon off;"]
