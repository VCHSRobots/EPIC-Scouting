# A simple bash script to deploy to git and the dev server.

# go fmt should be run by your IDE.

echo "[update.sh] Building."

echo "[update.sh] Latest build number is $(date +"%Y.%j.%H.%S")"

echo $(date +"%Y.%j.%H.%S") > .latest-build

env CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -v

# echo "[update.sh] Pushing to Git."

# git add .
# git commit
# git push

echo "[update.sh] Pushing to development server."

rsync -ahs --delete-before --info=progress2 --stats -e "ssh -p 2020" --exclude 'update.sh' --exclude 'config.yaml' --exclude 'campaigns.db' --exclude 'teams.db' --exclude 'users.db' --exclude '.latest-build' ".." "pi@pkre.co:/home/pi/src/"
rsync -ahs --delete-before -q -e "ssh -p 2020" --exclude 'update.sh' --exclude 'campaigns.db' --exclude 'teams.db' --exclude 'users.db' --exclude 'config.yaml' ".." "pi@pkre.co:/home/pi/src/" # Send .latest-build last.

echo "[update.sh] Pushing to production server."

rsync -ahs --delete-before --info=progress2 --stats -e "ssh" --exclude 'update.sh' --exclude 'config.yaml' --exclude 'campaigns.db' --exclude 'teams.db' --exclude 'users.db' --exclude '.latest-build' ".." "root@epicscouts.org:/root/go/src/"
rsync -ahs --delete-before -q -e "ssh" --exclude 'update.sh' --exclude 'config.yaml' --exclude 'campaigns.db' --exclude 'teams.db' --exclude 'users.db' ".." "root@epicscouts.org:/root/go/src/" # Send .latest-build last.

echo "[update.sh] Done!"
