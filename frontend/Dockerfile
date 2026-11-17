# ---------- Build ----------
FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

RUN npm run build


# ---------- Runtime ----------
FROM node:22-alpine

WORKDIR /app

RUN apk add --no-cache gettext && \
    npm install -g serve

COPY --from=builder /app/dist ./dist
COPY public/config.json.template ./config.json.template

EXPOSE 80

CMD sh -c "envsubst '\$DONATE_YOOMONEY_URL \$DONATE_CLOUDTIP_URL' < config.json.template > dist/config.json && exec serve -s dist -l tcp://0.0.0.0:80"
