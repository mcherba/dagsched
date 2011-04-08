
CC = 6g
LINK = 6l

CC32 = 8g
LINK32 = 8l

all dagsched: dagsched.go dag parser sorter
	$(CC) dagsched.go
	$(LINK) -o dagsched dagsched.6
	
dag: dag.go
	$(CC) dag.go
	
parser: parser.go
	$(CC) parser.go
	
sorter: sorter.go parser
	$(CC) sorter.go

clean:
	rm -f *.6
	rm -f dagsched

all32 dagsched32: dagsched.go dag parser32 sorter32
	$(CC32) dagsched.go
	$(LINK32) -o dagsched dagsched.8
	
dag32: dag.go
	$(CC32) dag.go
	
parser32: parser.go
	$(CC32) parser.go
	
sorter32: sorter.go parser32
	$(CC32) sorter.go
	
clean32:
	rm -f *.8
	rm -f dagsched
	
