#!/usr/bin/env sh

# abort on errors
set -e

# copy changes.md
cp ../changes.md ./zh/product/

# build
npm run docs:build

ossutil rm oss://docs-aiflowy-tech-website/ -rf
ossutil cp -rf assets/images oss://docs-aiflowy-tech-website/assets/images
ossutil cp -rf .vitepress/dist  oss://docs-aiflowy-tech-website/
