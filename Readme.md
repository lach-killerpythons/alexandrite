# ALEXANDRITE

working tools in Go for PSQL and REDIS

## Connect to PSQL DB

`	// set the credentials file location
  jade := JADE.JADE_FILE{"db.json", cwd + "/req"}
  // Return BLUE DB objct (psql)	
	DBx, err := BLUE.DB_JadeConnect("local", jade) `

## Connect to Redis DB

` // connect to example "localhost" pw:"" db:0 
	red := NewRedDB("localhost", "", 0)
	rdb, err := red.New_Connection() `


# BLUE
PSQL

DB struct 
- p *sql.DB
- credentials (Host, Port, User, DBName)
- Table

// CONNECT TO DB 
DB.Connect_DB(name) -> return DB with credentials

// CONNECT TO TABLE
DB.OpenTable(name)  -> connect to table

// CREATE A NEW TABLE
CREATE_TABLE_v1(name, []cols, []example_vals .. P_KEY, F_KEY)


Table struct
- Name, Cols

primary key
foreign key



# RED
REDIS DB

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
