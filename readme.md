A fairly simple reverse-proxy and cache server writting in GoLang, primarily inspired by [Varnish](https://varnish-cache.org/)

## Goals
- [x] Request coalescing
- [x] Vary support
- [ ] Respect cache control header
  - [x] private, public, no-cache, no-store, max-age, s-maxage
  - [ ] no-cache (properly this time), stale-if-error, stale-while-revalidate
- [ ] Query string sorting
- [ ] Latency-based origin loadbalancing
- [x] Cache tag/purge by tag support
- [ ] Memory mapped cache(?)

## Notes

- How do you properly do request coalescing for requests that can vary? e.g. if a bunch of requests come in for an image, they get coalesced into a single origin request, but that origin request may contain an `Accept` header that allows for a mime-type that not all requests that were coalesced together accept.
  - You can maybe do this somewhat intelligently by only coalescing requests together if they share accept and accept-encoding values?
- Implementing stale-if-error and stale-while-revalidate seem like they'll be a pain. You'd need to intelligently keep things in cache after they've expired

## Performance

Results of some benchmarks from requesting the same asset repeatedly (all cache hits)
| Request rate | Total requests made | Avg response time | Median response time | Min response time | Max response time | Response time Std. Dev. |
|--------------|---------------------|-------------------|----------------------|-------------------|-------------------|-------------------------|
| ~39k RPM     | 11167               | 86.1µs            | 69µs                 | 17µs              | 4.7ms             | 144.6µs                 |
| ~97k RPM     | 40899               | 63.1µs            | 37.8µs               | 15.8µs            | 9.5ms             | 134µs                   |
## Features

### Request Coalescing

If multiple requests for the same cache entry are received at the same time, only

### Full support for `Vary` header

Varies the cache based on values of headers listed in the `Vary` response header. There are currently no disallowed headers, even things like `User-Agent`
