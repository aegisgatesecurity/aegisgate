# Check what might cause exit code 1 in tests
# Run tests and check for any issues
import subprocess
import os

os.chdir('C:/Users/Administrator/Desktop/Testing/padlock')

# Run tests with verbose output
result = subprocess.run(['go', 'test', '-v', '-coverprofile=coverage.out', './...'], 
                       capture_output=True, text=True, timeout=300)

print("STDOUT length:", len(result.stdout))
print("STDERR length:", len(result.stderr))
print("Return code:", result.returncode)

# Check for FAIL in output
if 'FAIL' in result.stdout:
    print("\n=== FAIL found in output ===")
    lines = result.stdout.split('\n')
    for i, line in enumerate(lines):
        if 'FAIL' in line:
            # Print context around the FAIL
            start = max(0, i-2)
            end = min(len(lines), i+3)
            for j in range(start, end):
                print(lines[j])
            break

# Check stderr for errors
if result.stderr:
    print("\n=== STDERR ===")
    print(result.stderr[:2000])
