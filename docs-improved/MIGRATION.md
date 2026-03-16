# Documentation migration notes

This document explains what changed between `docs/` (old) and `docs-improved/`
(new), and why.

## Design principles

The rewrite follows patterns from well-maintained OSS projects (GitHub CLI,
Terraform, Cobra):

1. **README is a signpost, not an encyclopedia.** Keep it under 100 lines.
   Link to docs for details.
2. **Each document has one purpose and one audience.** Don't mix user guides
   with contributor references.
3. **No redundancy.** Say it once, link elsewhere. The old docs repeated CI/CD
   setup in README, installation.md, and ci-cd-guide.md.
4. **Factual tone.** Describe what the tool does. Avoid marketing language,
   benefit lists, and superlatives.
5. **Example-driven.** Show the command, then explain if needed. Don't explain
   before showing.

## Structure comparison

### Old (docs/)

| File | Lines | Issue |
|------|-------|-------|
| `README.md` | 635 | Tries to cover everything: quick start, CI/CD, architecture, marketing |
| `ci-cd-guide.md` | 706 | Duplicates README CI/CD section; marketing tone ("solves critical challenges") |
| `ENCRYPTED_STORAGE_ARCHITECTURE.md` | 880 | Internal design proposal, not user docs |
| `presentation.md` | 1250 | Slide deck, not documentation |
| `PATTERN_ROUTING.md` | 375 | Duplicates configuration.md |
| `GITCONFIG_COMMAND.md` | 497 | Single-command reference, could be in installation.md |
| `DIAGNOSTIC_LOGGING.md` | 529 | Could be a section in troubleshooting.md |
| `TOKEN_CACHING.md` | 305 | Internal details, relevant parts in architecture.md |
| `SLSA_IMPLEMENTATION_PLAN.md` | 387 | Implementation plan, not documentation |
| `JSON_MULTI_ORG_CONFIG.md` | 375 | Niche example, merged into configuration.md |
| `origin_of_the_project.md` | 57 | Pitch/marketing document |
| `TESTING.md` | 450 | Merged into development.md |
| `E2E_TESTING_TUTORIAL.md` | 396 | Internal test infrastructure |
| `E2E_INFRASTRUCTURE.md` | 270 | Internal test infrastructure |

**Total: ~6,700 lines across 14 docs + 635-line README**

### New (docs-improved/)

| File | Lines | Audience |
|------|-------|----------|
| `README.md` | ~90 | Everyone (entry point) |
| `installation.md` | ~100 | Users |
| `configuration.md` | ~180 | Users |
| `ci-cd.md` | ~150 | DevOps / CI users |
| `security.md` | ~80 | Users / security reviewers |
| `troubleshooting.md` | ~130 | Users |
| `architecture.md` | ~110 | Contributors |
| `development.md` | ~160 | Contributors |

**Total: ~1,000 lines across 8 files**

## What was removed and why

| Old file | Reason |
|----------|--------|
| `presentation.md` + `presentation-custom.css` | Slide deck, not docs. Keep separately if needed for talks. |
| `ENCRYPTED_STORAGE_ARCHITECTURE.md` | Design proposal / analysis. Move to a wiki or ADR folder if you want to preserve it. |
| `SLSA_IMPLEMENTATION_PLAN.md` | Implementation plan, not user-facing. Track in issues or a separate planning folder. |
| `origin_of_the_project.md` | Marketing/pitch content. The README's "How it works" section covers the same ground factually. |
| `JSON_MULTI_ORG_CONFIG.md` | Examples merged into `configuration.md`. |
| `PATTERN_ROUTING.md` | Merged into `configuration.md` (pattern matching section). |
| `GITCONFIG_COMMAND.md` | Merged into `installation.md` (gitconfig section). |
| `DIAGNOSTIC_LOGGING.md` | Merged into `troubleshooting.md` (debug logging section). |
| `TOKEN_CACHING.md` | Key facts in `architecture.md` (token lifecycle table). |
| `TESTING.md` | Merged into `development.md`. |
| `E2E_TESTING_TUTORIAL.md` | Internal. Reference from `development.md` if needed. |
| `E2E_INFRASTRUCTURE.md` | Internal test infra docs. Keep in `test/` if needed. |

## What was kept (merged)

All user-facing information from the old docs is preserved in the new structure.
Nothing was lost — it was consolidated and de-duplicated.

## Specific tone changes

### Before (old README)

```text
This extension is designed to solve common CI/CD authentication challenges
with GitHub Apps, including cross-organization access, git submodules, and
long-running jobs.

**Benefits:**
- No temporary files created
- No chmod needed
- No cleanup required
- Keys stored securely in memory only

### Advantages Over Robot Accounts
- GitHub Apps don't consume user licenses
- Fine-grained permissions per repository/organization
- Actions attributed to the GitHub App in audit logs
```

### After (new README)

```text
The extension registers itself as a Git credential helper. When Git needs
credentials for a URL, the helper matches the URL against configured patterns,
picks the credential with the longest matching prefix, and returns a
username/password pair.
```

The new version states what happens. The reader can decide whether it solves
their problem.

## Applying the new docs

To replace the old documentation:

1. Review the `docs-improved/` directory.
2. Copy files to their final locations:
   - `docs-improved/README.md` → root `README.md`
   - `docs-improved/*.md` (others) → `docs/`
3. Archive or delete old files that were not carried forward.
4. Update any internal links.
5. Run `make markdownlint` to verify formatting.
