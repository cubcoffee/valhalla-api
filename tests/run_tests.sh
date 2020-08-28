#!/bin/bash
echo "############   Iniciando Testes   ############"
cd ..
go test -v -coverprofile=cover.out -coverpkg=github.com/cubcoffee/valhalla-api,github.com/cubcoffee/valhalla-api/model,github.com/cubcoffee/valhalla-api/router,github.com/cubcoffee/valhalla-api/dao ./...
go tool cover -html=cover.out -o cover.html
rm cover.out