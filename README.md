# gostmplcli
A simple go template command line with zero dependencies

## Install
go get -u github.com/jesusnoseq/gostmplcli

## Example to execute from binary file
```shell script
./gostmplcli -r template_c.input -t test_data/template_a.input -t test_data/template_b.input -t test_data/template_c.input
```

## Example to execute from Docker image
```shell script
docker run --rm -v ${pwd}:/app -it gostmplcli -r template_c.input -t test_data/*.input
```
