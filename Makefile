MODULE   := $(shell awk 'NR==1{print $$2}' go.mod)
PIGEON   := go run -mod=vendor github.com/mna/pigeon
DIST     := dist
BIN_NAME := $(notdir $(MODULE))
BIN_PATH := $(DIST)/$(BIN_NAME) 
FLAGS    := -mod=vendor

grammar: svg/ast/dattr_peg.go

svg/ast/dattr_peg.go: svg/ast/dattr.peg svg/ast/ast.go
	$(PIGEON) -o $@ svg/ast/dattr.peg

.PHONY: $(BIN)  # Let `go` use its own caching
$(BIN):
	go build $(FLAGS) -o $(BIN)
