#!/bin/bash
cd ../at-admin;git add .;git commit -a -m up;git pull;git push --all;
cd ../at-api;git add .;git commit -a -m up;git pull;git push --all;
cd ../at-client;git add .;git commit -a -m up;git pull;git push --all;
cd ../at-cluster;git add .;git commit -a -m up;git pull;git push --all;
cd ../at-common;git add .;git commit -a -m up;git pull;git push --all;
cd ../at-docs;git add .;git commit -a -m up;git pull;git push --all;
cd ../at-server;git add .;git commit -a -m up;git pull;git push --all;
cd ../at-web;git add .;git commit -a -m up;git pull;git push --all;
