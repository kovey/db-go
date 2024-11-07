## ksql tool
#### Description
###### The tool can process database changed between dev and prod
###### Install
    go install github.com/kovey/db-go/ksql@latest
### Examples
    ksql -h
#### diff table changed
    ksql diff --dir=path/to/sql --driver=mysql --todb=test --to='root:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4' --fromdb=test_prod --from='root:password@tcp(127.0.0.1:3306)/test_prod?charset=utf8mb4'
