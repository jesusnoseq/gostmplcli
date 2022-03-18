[![Lint](https://github.com/jesusnoseq/gostmplcli/actions/workflows/lint.yml/badge.svg)](https://github.com/jesusnoseq/gostmplcli/actions/workflows/lint.yml)
[![Unit test](https://github.com/jesusnoseq/gostmplcli/actions/workflows/test.yml/badge.svg)](https://github.com/jesusnoseq/gostmplcli/actions/workflows/test.yml)


# gostmplcli
A simple go template command line with zero dependencies

## Install
```shell script
go get -u github.com/jesusnoseq/gostmplcli
```

## Example from binary file
```shell script
./gostmplcli -r template_c.input -t test_data/template_a.input -t test_data/template_b.input -t test_data/template_c.input
```

## Example from Docker image
```shell script
docker run --rm -v ${pwd}:/app -it gostmplcli -r template_c.input -t test_data/*.input
```
