#!/bin/sh

baseURL=https://your-example.com/
logfile=/var/log/blog-updater.log

cd /opt/personal-blog
curDirName=$(date +"%Y%m%dT%H%M%S")
destDir=public/all/$curDirName

echo ---- $(date) ---- >> $logfile

git pull >> $logfile \
 && git submodule update >> $logfile \
 && hugo --cleanDestinationDir --baseURL $baseURL --destination $destDir >> $logfile

if [ $? -eq 0 ]; then
    cd /opt/personal-blog/public
    ln -nfs ./all/$curDirName current

    cd /opt/personal-blog/public/all
    find . -mindepth 1 -name $curDirName -type d -prune -o -exec rm -rf {} +
fi
