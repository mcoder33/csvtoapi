Sending data from CSV to api

-api string
    Target API endpoint
-bs int
    Buffer size for line (default 1024)
-colMap string
    Field colMap like key:value,key2:value2
-debugMode
    Debug mode instead sending real requests just log
-file string
    Put here the full file path
-rps int
    RPS - good speed for you (default 10)
-sep string
    Separator for words in line (default ",")

Example command:
csvtoapi-mac-arm64 -api https://youre.sit/some/action -file ./tmp/test.csv -colMap someData:someId,customId:newId,anyId:externalAnyId,someStatus:status,createDate:attributionDate,sumOther:sum
