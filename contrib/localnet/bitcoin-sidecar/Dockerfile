FROM node:18.20.4 as builder

WORKDIR /home/zeta/node

COPY bitcoin-sidecar/js/* .

RUN npm install && npm install typescript -g && tsc

FROM node:alpine

COPY --from=builder /home/zeta/node/dist ./dist
COPY --from=builder /home/zeta/node/node_modules ./node_modules

CMD ["node", "dist/index.js"]