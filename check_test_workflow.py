# Check test.yml workflow
with open('.github/workflows/test.yml', 'r') as f:
    content = f.read()
    print(content)
