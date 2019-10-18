#!/bin/sh

# System startup

baseURL=https://your-example.com/
logfile=/var/log/blog-updater.log

cd /opt/personal-blog
curDirName=$(date +"%Y%m%dT%H%M%S")
destDir=/var/cache/nginx/blog/$curDirName

echo ---- $(date) ---- >> $logfile

hugo --cleanDestinationDir --baseURL $baseURL --destination $destDir >> $logfile

if [ $? -eq 0 ]; then
    # Soft link on SSD driver only
    cd /opt/personal-blog
    ln -nfs $destDir public

    # Remove old files
    cd $destDir/..
    find . -mindepth 1 -name $curDirName -type d -prune -o -exec rm -rf {} +
fi
