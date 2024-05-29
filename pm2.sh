go build
pm2 restart pm2.json && pm2 save && pm2 log
