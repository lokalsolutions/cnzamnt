FROM node:22-bookworm AS deps

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm install

FROM deps AS dev

COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

FROM deps AS build

COPY . .
RUN npm run build

FROM caddy:2-alpine AS prod

COPY --from=build /app/dist /usr/share/caddy

EXPOSE 80
