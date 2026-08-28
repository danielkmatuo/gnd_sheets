#!/usr/bin/env fish

# Start with a completely clean distribution directory
echo "Starting build process..."
rm -rf dist
echo "Removed dist dir..."

# Create release directories
echo "Creating binaries and data dirs..."
mkdir -p dist/gnd-sheets-linux-amd64/data
mkdir -p dist/gnd-sheets-windows-amd64/data

# Copy the current frontend into both distributions
echo "Copying frontend dir to binaries dirs..."
cp -r frontend dist/gnd-sheets-linux-amd64/frontend
cp -r frontend dist/gnd-sheets-windows-amd64/frontend

# Copy the current reference data into both distributions
echo "Copying reference data dir to binaries dirs..."
cp -r data/reference dist/gnd-sheets-linux-amd64/data/reference
cp -r data/reference dist/gnd-sheets-windows-amd64/data/reference

echo "Building G&D Sheets for linux as gnd_sheets..."
env GOOS=linux GOARCH=amd64 go build -o dist/gnd-sheets-linux-amd64/gnd_sheets .

if test $status -ne 0
    echo "Linux build failed"
    exit 1
end

echo "Building G&D Sheets for windows as gnd_sheets.exe..."
env GOOS=windows GOARCH=amd64 go build -o dist/gnd-sheets-windows-amd64/gnd_sheets.exe .

if test $status -ne 0
    echo "Windows build failed"
    exit 1
end

echo "Building .zip and .tar.gz for windows and linux dirs..."
cd dist

zip -r gnd-sheets-windows-amd64.zip gnd-sheets-windows-amd64
tar -czf gnd-sheets-linux-amd64.tar.gz gnd-sheets-linux-amd64

cd ..

echo "Build complete!"
