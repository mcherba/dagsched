
CC = 6g
LINK = 6l

CC32 = 6g
LINK32 = 6l

all dagsched: dagsched.go dag
	$(CC) dagsched.go
	$(LINK) -o dagsched dagsched.6
	
dag: dag.go
	$(CC) dag.go

clean:
	rm -f *.6
	rm -f dagsched

all32 dagsched32: dagsched.go dag
	$(CC32) dagsched.go
	$(LINK32) -o dagsched dagsched.8
	
dag32: dag.go
	$(CC32) dag.go

clean32:
	rm -f *.8
	rm -f dagsched
	
