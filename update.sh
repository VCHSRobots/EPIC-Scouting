# A simple bash script to deploy to git and the dev server.

# go fmt should be run by your IDE.

echo "[update.sh] Building."

echo "[update.sh] Latest build number is"

echo $(date +"%F %T.%N") > .latest-build

env CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -v

# echo "[update.sh] Pushing to Git."

# git add .
# git commit
# git push

echo "[update.sh] Pushing to development server."

rsync -ahs --delete-before --info=progress2 --stats -e "ssh -p 2020" --exclude 'update.sh' --exclude 'config.yaml' ".." "pi@pkre.co:/home/pi/src/"

echo "[update.sh] Done!"
