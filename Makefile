
CC = 6g
LINK = 6l


all dagsched: dagsched.go
	$(CC) dagsched.go
	$(LINK) -o dagsched dagsched.6

clean:
	rm -f *.6
	rm -f dagsched
	
