# Developer Certificate of Origin (DCO)

> **⚠️ IMPORTANT: All commits must include a Signed-off-by line.**
> 
> This follows the Linux Foundation DCO 1.1 standard.

---

## Version 1.1

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I have the right to submit it under the open source license indicated in the file; or

(b) The contribution is based upon previous work that, to the best of my knowledge, is covered under an appropriate open source license and I have the right under that license to submit that work with modifications, whether created in whole or in part by me, under the same open source license (unless I am permitted to submit under a different license), as indicated in the file; or

(c) The contribution was provided directly to me by some other person who certified (a), (b) or (c) and I have not modified it.

(d) I understand and agree that this project and the contribution are public and that a record of the contribution (including all personal information I submit with it, including my sign-off) is maintained indefinitely and may be redistributed consistent with this project or the open source license(s) involved.

---

## Sign-Off Requirement

**Every commit must include a Signed-off-by line:**

```
Signed-off-by: Your Name <your.email@example.com>
```

### Using the -s Flag

Git automatically adds the Signed-off-by line when you use the `-s` flag:

```bash
# Good: Commit with sign-off
git commit -s -m "feat: add new feature"

# Bad: Missing sign-off (will be rejected)
git commit -m "feat: add new feature"
```

### Amending Missing Sign-Offs

If you forget to add sign-off to a commit:

```bash
# Amend the last commit to add sign-off
git commit --amend --signoff

# Amend without changing message
git commit --amend --no-edit --signoff
```

---

## Why DCO?

1. **Legal Protection**: Ensures contributors have the right to contribute
2. **Chain of Trust**: Creates a traceable lineage for code
3. **Open Source Standard**: Used by Linux Kernel, Docker, Kubernetes, and 1000+ projects

---

## Enforced By

- **GitHub Actions**: Automated checks reject unsigned commits
- **Maintainer Review**: PRs without sign-off will not be merged

---

## Contact

For questions about the DCO, contact: **legal@aegisgate.io**

---

**By contributing to AegisGate, you agree to sign off on all commits.**