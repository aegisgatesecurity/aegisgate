# Create commit
import subprocess
import os

os.chdir('C:/Users/Administrator/Desktop/Testing/padlock')

# Stage the files
subprocess.run(['git', 'add', 'cmd/padlock/main.go', 'pkg/core/core_test.go'], check=True)

# Get the staged diff
result = subprocess.run(['git', 'diff', '--cached', '--stat'], capture_output=True, text=True)
print("Staged files:", result.stdout)

# Create commit using subprocess with proper message
result = subprocess.run(['git', 'commit', '-m', 'fix: version format and test fixes for CI'], 
                       capture_output=True, text=True)
print("STDOUT:", result.stdout)
print("STDERR:", result.stderr)
print("Return code:", result.returncode)
