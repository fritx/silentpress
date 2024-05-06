# Prepare
cp -r p_example p
cp .env.example .env
# set your own config & secrets
# vim .env
# for example
sed -i.bak "s|^COOKIE_SECRET=.*|COOKIE_SECRET=\"$(openssl rand -base64 32)\"|" .env
sed -i.bak 's|^ADMIN_PASSWORD=.*|ADMIN_PASSWORD=\"SilentPress\"|' .env
sed -i.bak 's/^PORT=.*/PORT=8082/' .env

# Install dependencies
(cd silent && git stash -u)
git submodule update --init --recursive
(cd silent && git apply ../silent.patch)
go mod download

# Add a cron job to recover ./p contents
# How to create a cron job using Bash automatically without the interactive editor?
# https://stackoverflow.com/questions/878600/how-to-create-a-cron-job-using-bash-automatically-without-the-interactive-editor
# write out current crontab
crontab -l > mycron
# echo new cron into cron file
echo "*/3 * * * * cd ~/m/i/silentpress && cp -r p_example/** p/" >> mycron
# install new cron file
crontab mycron
rm mycron
# Cron line explaination
# * * * * * "command to be executed"
# - - - - -
# | | | | |
# | | | | ----- Day of week (0 - 7) (Sunday=0 or 7)
# | | | ------- Month (1 - 12)
# | | --------- Day of month (1 - 31)
# | ----------- Hour (0 - 23)
# ------------- Minute (0 - 59)

# Start or restart later..
gspp
(cd silent && git stash -u)
git submodule update --init --recursive
(cd silent && git apply ../silent.patch)
pm2 start pm2.json && pm2 log
