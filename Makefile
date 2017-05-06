BIN_TARGETS = depgrapher $(TRACEE)
TEST_TARGETS = $(TESTCMD) $(SEGFAULT)
TEST_TEMPS = a b c
OBJECT_FILES = */*.o

TRACEE = tracee/tracee
TESTCMD = testcmd/testcmd
SEGFAULT = testcmd/segfault

default: compile

clean:
	rm -f $(BIN_TARGETS) $(TEST_TARGETS) $(TEST_TEMPS) $(OBJECT_FILES)

compile: $(BIN_TARGETS)

test: $(BIN_TARGETS) $(TEST_TARGETS)
	./depgrapher -tracee $(TRACEE) -d /dev/stdout $(TESTCMD)
	./depgrapher -tracee $(TRACEE) -t /dev/stdout $(SEGFAULT)

install: $(BIN_TARGETS)
	mkdir -p $(DESTDIR)
	cp $(BIN_TARGETS) $(DESTDIR)

depgrapher: *.go
	go build -o $@
