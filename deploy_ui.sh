#!/bin/bash

# Deploy the UI
cd frontend
npm run build
cd ..

# Copy the UI to the web directory
echo "Copying the UI to the web directory..."
rm -rf goli/web/assets/*
cp -r web/* goli/web/