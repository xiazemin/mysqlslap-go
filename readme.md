a tool like mysqlslap implement in golang

useage:
```
 mysqlslap  -Hhost -uuser -ppassword -P3306 -ddatabase -q"select id from deviceattr where name='attr10' or name='attr20' group by id;" -ffilename -c 50 -i 100
 ```

 输出报告：
 ```
 Summary:
  Name:         mysqlslap
  Count:        40
  Total:        1.01 s
  Slowest:      1.00 s
  Fastest:      1.01 s
  Average:      1.00 s
  Requests/sec: 39.76

Response time histogram(ms):
  1003.000 [1]  |∎∎
  1003.300 [0]  |
  1003.600 [0]  |
  1003.900 [0]  |
  1004.200 [24] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  1004.500 [0]  |
  1004.800 [0]  |
  1005.100 [12] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  1005.400 [0]  |
  1005.700 [0]  |
  1006.000 [3]  |∎∎∎∎∎

Latency distribution:
  10 % in 1.00 s 
  25 % in 1.00 s 
  50 % in 1.00 s 
  75 % in 1.00 s 
  90 % in 1.00 s 
  95 % in 1.01 s 
  99 % in 1.01 s 

Status code distribution:
  [success]   40 responses   
  [failed]    0 responses 
 ```