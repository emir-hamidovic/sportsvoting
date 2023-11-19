FROM node:17-alpine AS development

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . ./

ENV PORT=3000

CMD ["npm", "start"]

FROM development AS build

RUN npm run build

FROM nginx:alpine

COPY --from=build /app/.nginx/nginx.conf /etc/nginx/conf.d/default.conf

WORKDIR /usr/share/nginx/html

RUN rm -rf ./*

COPY --from=build /app/build .

ENTRYPOINT ["nginx", "-g", "daemon off;"]