# Xerophi
The future of the CactusAPI

## Branches
All features being developed on should be worked on on branches with the naming
format `feature/[short but descriptive name]`, bugfixes with the format
`bug/[bug name (- issue # here if applicable)]`, and hot patches with the format
`hotfix/[short bug description/issue #]`

### Specific branches
- `develop`
 
  This branch is used for any semi-stable code, that is working, but not
  necessarily production-ready. Avoid committing code directly to `develop`,
  but it is allowed.

- `master`

  Stable, production-ready code only. Code may not be committed directly to
  `master`, it must be created on another branch and a merge/pull request
  created.



## Breaking Changes
Since this is a major rewrite of the API there are some changes that have been
made that are breaking to v1 compatibility. As much was kept the same possible
but some things weren't possible to implement (at least initially) and others
were re-thought. This API will be available under /api/v2/.

- `?append=true` parameter on config editing no longer works
    
  Any appends/removals/edits need to retrieve the current version from the DB and then add/remove from that as needed