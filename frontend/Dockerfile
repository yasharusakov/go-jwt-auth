FROM node:alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

ARG VITE_API_URL
ARG APP_PORT

RUN npm run build

FROM nginx:alpine

RUN apk --no-cache add curl

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE ${APP_PORT}

CMD ["nginx", "-g", "daemon off;"]