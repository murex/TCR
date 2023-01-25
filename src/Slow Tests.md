# Slow Tests

## Capturing test duration

```shell
cd .../src
go test -count=1 -tags=test_helper -v -json ./... > testoutput.txt
```

## Converting to CSV

```shell
cat testoutput.txt | \
  jq -r 'select(.Action == "pass" and .Test != null) | .Test + "," + (.Elapsed | tostring)' | \
   sort --reverse -k2 -n -t, | tee testoutput.csv
```

## Import into excel
