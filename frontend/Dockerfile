# ---------- Build ----------
FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

ARG VITE_DONATE_YOOMONEY_URL
ARG VITE_DONATE_CLOUDTIP_URL

ENV VITE_DONATE_YOOMONEY_URL=$VITE_DONATE_YOOMONEY_URL
ENV VITE_DONATE_CLOUDTIP_URL=$VITE_DONATE_CLOUDTIP_URL

RUN npm run build


# ---------- Runtime ----------
FROM node:22-alpine

WORKDIR /app

RUN npm install -g serve

COPY --from=builder /app/dist ./dist

EXPOSE 80

CMD ["serve", "-s", "dist", "-l", "tcp://0.0.0.0:80"]