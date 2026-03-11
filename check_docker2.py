# Check exact docker section in test.yml
with open('.github/workflows/test.yml', 'r') as f:
    lines = f.readlines()

for i in range(69, len(lines)):
    print(f"{i+1}: {repr(lines[i])}")
