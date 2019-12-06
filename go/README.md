# Brindexer


# Brindexer is a Golang based index engine that scans file system and index files into sqlite DBs

1. Requirement: GO 1.12.X (or above)

2. Build:
   cd brindex/go (whereever you file tree is)
   ./prepare (you only need to this once)
   ./build


3. To test

	 cd brindex/go ( top of where your go src tree is.)

    
	# Scan and indexing:


	bin/index  file_directory_to_index



	# Parallel scan and indexing: (we will have separate doc for this)


	bin/parallelindex  file_directory_to_index


	# After index completes, the following gives summary of the files indexed:

	bin/summary  indexed-dir

	# Run query:

	bin/query -q "select %s from %s where size > 100000" indexed-dir

	Note: use %s for general purpose queries. If you know the DB schema, you can use -a option to run SQL statements as is. (e.g ./query -a -i  "analyze" /indexed-dir) to perform sqlite optimization for all DBs





	Build a docker image:

	cd brindexer/go ( top of where your src tree is.)

	docker-compose build
