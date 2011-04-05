
CC = 8g
LINK = 8l


all dagsched: dagsched.go
	$(CC) dagsched.go
	$(LINK) -o dagsched dagsched.8

clean:
	rm -f *.8
	rm -f dagsched
	
