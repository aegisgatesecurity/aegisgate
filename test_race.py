# Run tests with CGO enabled
import subprocess
import os

os.chdir('C:/Users/Administrator/Desktop/Testing/padlock')

# Set CGO_ENABLED=1 and run tests
env = os.environ.copy()
env['CGO_ENABLED'] = '1'

# Run tests similar to CI
result = subprocess.run(['go', 'test', '-v', '-race', '-coverprofile=coverage.out', './pkg/core/...'], 
                       capture_output=True, text=True, env=env, timeout=120)

print("Return code:", result.returncode)
print("STDOUT (last 2000 chars):", result.stdout[-2000:])
print("STDERR:", result.stderr[:500] if result.stderr else "(none)")
