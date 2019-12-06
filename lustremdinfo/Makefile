all: lustremdinfo test

lustremdinfo: lustremdinfo.o lustremdinfo.h
	cc -c lustremdinfo.c -llustreapi
	ar rsv lustremdinfo.a lustremdinfo.o

test: lustremdinfo.o test.o lustremdinfo.h
	cc -o test lustremdinfo.c test.c -llustreapi

clean:
	rm -f lustremdinfo lustremdinfo.o lustremdinfo.a test test.o

.PHONY: all test clean
