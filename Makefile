build:
	for i in day* ; do ( cd $$i; echo "--> $$i"; go build ); done

clean:
	for i in day* ; do ( cd $$i; echo "--> $$i"; go clean ); done

# eg. make run INPUT=example01.txt
INPUT := input.txt
run:
	for i in day*; do ( cd $$i; echo "--> $$i ${INPUT}"; ./$$i ${INPUT} ); done
