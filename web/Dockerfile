FROM node:lts-jod AS builder
RUN npm install -g pnpm@latest-10
WORKDIR /app
COPY . .
RUN pnpm install --frozen-lockfile
RUN pnpm run build
RUN pnpm prune --prod # prune non-dev dependencies

FROM node:jod-alpine
WORKDIR /app
EXPOSE 4173
CMD [ "node", "index.js" ]
ENV NODE_ENV=production
ENV ORIGIN=http://localhost:4173
ENV PORT=4173
ARG ACCESS_TOKEN
ENV ACCESS_TOKEN=${ACCESS_TOKEN}
ARG API_BASE_URL
ENV API_BASE_URL=${API_BASE_URL}
COPY --from=builder /app/package.json .
COPY --from=builder /app/build .
COPY --from=builder /app/node_modules node_modules
