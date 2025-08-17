# [Arg variable]
ARG WORK_DIR="/app"
ARG APP_DIR="./apps/front"

FROM node:20-alpine AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

# build (for development)
FROM base AS develop
ARG WORK_DIR
ARG APP_DIR
WORKDIR ${WORK_DIR}

COPY ${APP_DIR}/package.json ./
COPY ${APP_DIR}/pnpm*.yaml ./

RUN pnpm install --frozen-lockfile

COPY ${APP_DIR} .

FROM develop AS builder
ARG WORK_DIR
WORKDIR ${WORK_DIR}

RUN pnpm build

# release
FROM alpine:latest AS app
ARG WORK_DIR
ENTRYPOINT ["app"]
