#!/bin/sh

baseURL=https://your-example.com/
logfile=/var/log/blog-updater.log

cd /opt/personal-blog

echo ---- $(date) ---- >> $logfile

git pull >> $logfile \
 && git submodule update >> $logfile \
 && hugo --cleanDestinationDir --baseURL $baseURL  >> $logfile
