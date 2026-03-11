# Check ci.yml for test configuration
with open('.github/workflows/ci.yml', 'r') as f:
    content = f.read()
    print(content)
