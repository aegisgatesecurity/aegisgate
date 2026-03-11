# Create commit - only core_test.go since main.go is ignored
import subprocess
import os

os.chdir('C:/Users/Administrator/Desktop/Testing/padlock')

# Stage only the core_test.go
subprocess.run(['git', 'add', 'pkg/core/core_test.go'], check=True)

# Create commit using subprocess with proper message
result = subprocess.run(['git', 'commit', '-m', 'fix: test fixes for CI - use TierCommunity'], 
                       capture_output=True, text=True)
print("STDOUT:", result.stdout)
print("STDERR:", result.stderr)
print("Return code:", result.returncode)
