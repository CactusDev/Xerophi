# Xerophi
The future of the CactusAPI

## Breaking Changes
Since this is a major rewrite of the API there are some changes that have been
made that are breaking to v1 compatibility. As much was kept the same possible
but some things weren't possible to implement (at least initially) and others
were re-thought. This API will be available under /api/v2/.

- `?append=true` parameter on config editing no longer works
    
  Any appends/removals/edits need to retrieve the current version from the DB and then add/remove from that as needed