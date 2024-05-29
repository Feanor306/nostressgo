# NOSTRess Go
NOSTRess Go is a [NOSTR](https://nostr.com/) relay written in Go
Using websockets, postgresql and pgadmin

# NIPs implemented
- [x] 1 => basic protocol
- [x] 12,16,20,33 => merged into 1
- [x] 9 => event deletion 
- [x] 11 => relay information document 
- [x] 14 => event subject
- [x] 40 => event expiration
- [x] 50 => search

# Data permissions for dockerized volumes
```sh
# postgres
# env var PGDATA is important to avoid permission issues when mounting
PGDATA=/var/lib/postgresql/data/pgdata
# can still mount /var/lib/postgresql/data
# ./data/pgdata:/var/lib/postgresql/data

# pgadmin permissions
sudo chown -R data/pgadmin 5050:5050
```

## Potential NIP-s to implement
* 2 => follow list (for backup purposes)
* 4 => encrypted direct message
* 8 => mentions
* 15 => marketplace
* 18 => reposts
* 23 => long form content
* 24 => more metadata
* 25 => reactions
* 26 => event delegation (posting with another key)
* 27 => mentions in text
* 28 => public chat
* 32 => labeling / reporting
* 38 => user statuses
* 39 => user social media profiles
* 45 => event counts
* 56 => report