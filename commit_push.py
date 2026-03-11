import subprocess
import os

os.chdir('C:/Users/Administrator/Desktop/Testing/padlock')

# Stage Dockerfile
subprocess.run(['git', 'add', 'Dockerfile'])

# Commit
result = subprocess.run(['git', 'commit', '-m', 'fix: add go.sum to Docker build for dependencies'], 
                       capture_output=True, text=True)
print("Commit result:", result.returncode)
print("STDOUT:", result.stdout)
print("STDERR:", result.stderr)

# Push
if result.returncode == 0:
    print("\nPushing...")
    result = subprocess.run(['git', 'push'], capture_output=True, text=True)
    print("Push result:", result.returncode)
    print("STDOUT:", result.stdout)
    print("STDERR:", result.stderr)
