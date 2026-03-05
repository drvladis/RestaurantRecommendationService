package app

import (
	dl "APIforElasticBD/internal/db/dataloader"
	ic "APIforElasticBD/internal/db/indexcreator"
)

func ESCreateIndices(idxname string) {
	ic.IndexCreator(idxname)
}

func ESDataLoad(filename, index string) {
	dl.DataLoader(filename, index)
}
