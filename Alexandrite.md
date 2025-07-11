## ALEXANDRITE

# RED
REDIS DB

# BLUE
PSQL

# JADE
manage credentials for different DBs


// caching 

GET request
↓
Check Redis → HIT → Return
         ↓
       MISS
         ↓
Query Postgres → Save in Redis → Return
