#!/bin/sh
cd /opt/personal-blog
git pull && git submodule update && hugo
