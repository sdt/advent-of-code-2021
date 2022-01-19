build:
	for i in day* ; do ( cd $$i; echo "--> $$i"; go build ); done

clean:
	for i in day* ; do ( cd $$i; echo "--> $$i"; go clean ); done

fmt:
	for i in day* ; do ( cd $$i; echo "--> $$i"; go fmt ); done

vet:
	for i in day* ; do ( cd $$i; echo "--> $$i"; go vet ); done

# eg. make run INPUT=example01.txt
INPUT := input.txt
run:
	for i in day*; do ( cd $$i; echo "--> $$i ${INPUT}"; ./$$i ${INPUT} ); done

time:
	for i in day*; do ( cd $$i; echo "--> $$i ${INPUT}"; time ./$$i ${INPUT} 2>&1 ); done 2>&1 | egrep 'day|real'
