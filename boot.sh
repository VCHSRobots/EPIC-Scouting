# Call this script with CRON.

echo "Ping!"

# IF the program isn't running, try starting it.
running=$(ps aux | grep '[E]PIC-Scouting' | grep -v 'boot.sh')
if [[ -z $running ]]; then
    echo "Not running, therefore starting."
    go build . && ./EPIC-Scouting
fi
# Check for modified files.
latestBuild=$(<.latest-build)
lastSeenBuild=$(<.last-seen-build)
if [[ "$latestBuild" != "$lastSeenBuild" ]]; then
    echo "Detected update from $lastSeenBuild to $latestBuild. Restarting."
    sleep 1
    kill $(ps aux | grep '[E]PIC-Scouting' | grep -v 'boot.sh' | awk '{print $2}') # Shut down.
    sleep 1
    go build . && ./EPIC-Scouting
fi
# Update most recently seen build number.
cp .latest-build .last-seen-build
