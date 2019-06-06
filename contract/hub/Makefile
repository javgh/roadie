roadie:
	$(eval TMPDIR=$(shell mktemp -d))
	solc Hub.sol --bin --abi --optimize -o $(TMPDIR)
	cp $(TMPDIR)/Hub.abi .
	cp $(TMPDIR)/Hub.bin .
	abigen --bin=Hub.bin --abi=Hub.abi --pkg=hub --out=hub.go
