# Read binary to check raw bytes
with open('.github/workflows/release.yml', 'rb') as f:
    content = f.read()

lines = content.split(b'\n')
print(b"Line 63: " + lines[62])
