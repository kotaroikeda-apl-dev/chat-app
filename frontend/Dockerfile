# Stage 1: Build the app
FROM node:18 as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .

# ビルド用の環境変数を設定
ARG REACT_APP_URL_DOMAIN
ENV REACT_APP_URL_DOMAIN=$REACT_APP_URL_DOMAIN

RUN npm run build

# Stage 2: Serve the app
FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]