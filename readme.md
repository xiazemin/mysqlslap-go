a tool like mysqlslap implement in golang

useage:
```
 mysqlslap  -Hhost -uuser -ppassword -P3306 -ddatabase -q"select id from deviceattr where name='attr10' or name='attr20' group by id;" -ffilename -c 50 -i 100
 ```