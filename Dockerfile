FROM denoland/deno
WORKDIR /root
COPY . .
EXPOSE 80 443 9000 9001 9002
CMD ["deno","run","--allow-all",".\main.ts"]
