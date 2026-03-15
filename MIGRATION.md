# License Migration Guide v1.0.3

## Overview

Migrate from previous AegisGate versions to v1.0.3 with cryptographically signed licenses.

## Breaking Changes

**Old license keys are NOT VALID**

- **OLD**: Client-side validation with DEV_MODE bypass
- **NEW**: Cryptographic validation with HMAC/RSA signatures

## Migration Steps

### Step 1: Backup Configuration

```bash
# Backup environment
env | grep -i aegisgate > aegisgate-env-backup.txt

# Backup configs
cp config/aegisgate.yml config/aegisgate.yml.backup
cp .env .env.backup
```

### Step 2: Obtain New License

Contact administrator or sales with:
- Current license key (for lookup)
- Email address
- Tier (Developer/Professional/Enterprise)
- Hardware ID (Enterprise only)

### Step 3: Install v1.0.3

```bash
# Pull latest
git pull origin main

# Or download binary
curl -L https://github.com/aegisgatesecurity/aegisgate/releases/download/v1.0.3/aegisgate_v1.0.3.tar.gz | tar xz
```

### Step 4: Update License

Replace old key with new signed key:

```bash
# Old (example)
export AEGISGATE_LICENSE_KEY="pro-xxxxxxxx"

# New (example)
export AEGISGATE_LICENSE_KEY="base64(JSON).base64(SIGNATURE)"
```

### Step 5: Verify

```bash
# Check version
./aegisgate --version

# Verify license
./aegisgate --check-license

# Start
./aegisgate
```

## Troubleshooting

### "Invalid license signature"
- Copy complete license key
- Check expiry date
- Verify HMAC secret matches

### "Hardware ID mismatch" (Enterprise)
- License for different hardware
- Re-issue with correct hardware ID

### "License expired"
- Contact sales for renewal

## Support

- Email: support@aegisgate.io
- Subject: Migration v1.0.3 - [Organization]
