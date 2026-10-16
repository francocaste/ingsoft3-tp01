# ---- Etapa 1: build ----
FROM node:22-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# ---- Etapa de tests: se construye con --target test ----
# Reusa todo lo de "build" (dependencias + codigo) y al ejecutarse corre los
# tests con cobertura. Si falla un test o el umbral, el contenedor sale con
# codigo 1 y el job de CI falla. No entra en la imagen final.
FROM build AS test
ENTRYPOINT ["npm", "run", "test:ci"]

# ---- Etapa final: nginx sirve los estáticos ----
FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80