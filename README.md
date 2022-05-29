# GoFuzzy

Fuzzy search library backed by an in-memory [Trie](https://en.wikipedia.org/wiki/Trie). Memory is getting cheaper and
larger, reference datasets can be loaded completely in memory on servers and used both for:

1. Validating input data
2. Providing an autocomplete endpoint

```mermaid
flowchart LR
  A[User] -- request--> API
  subgraph Backend
    direction LR
    subgraph API
        direction TB
        ds[In Memory Dataset\nFor Validation\nand\nAutocomplete]
    end
    API --> id1[(Database)]
  end
```

The main use case of this library is to help on the autocomplete endpoint.
