.POSIX:

all:
	@make clean
	@cd cmd/indexcreator && go build -o indexcreator
	@cd cmd/indexcreator && ./indexcreator
	@cd cmd/indexcreator && rm indexcreator
	@cd cmd/dataloader && go build -o dataloader
	@cd cmd/dataloader && ./dataloader
	@cd cmd/dataloader && rm dataloader
	@cd cmd/http && go build -o http
	@cd cmd/http && ./http
	
indexcreator:
	@cd cmd/indexcreator && go build -o indexcreatorUtil
	@cd cmd/indexcreator && ./indexcreatorUtil
	@cd cmd/indexcreator && rm indexcreatorUtil

dataloader:
	@cd cmd/dataloader && go build -o dataloaderUtil
	@cd cmd/dataloader && ./dataloaderUtil
	@cd cmd/dataloader && rm dataloaderUtil

http:
	@cd cmd/http && go build -o httpUtil
	@cd cmd/http && ./httpUtil

clean:
	@-cd cmd/indexcreator && rm main indexcreator indexcreatorUtil
	@-cd cmd/dataloader && rm main dataloader dataloaderUtil
	@-cd cmd/http && rm main http httpUtil

.PHONY: all indexcreator dataloader http clean