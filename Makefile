GO_PROJECT_TEMP = $(HOME)/go_project_temp

all: change_gopath install  

########################################
### Build

install: 
	export GOPATH=$(GO_PROJECT_TEMP) && cd "$(GO_PROJECT_TEMP)/src/dbachain" && go install  ./cmd/dbachaincli
	export GOPATH=$(GO_PROJECT_TEMP) && cd "$(GO_PROJECT_TEMP)/src/dbachain" && go install  ./cmd/dbachaind


########################################
### Dependencies

change_gopath:
	@mkdir -p $(GO_PROJECT_TEMP)/src
	@ln -sf  $(GOPATH)/src/dbachain/ $(GO_PROJECT_TEMP)/src
