0. ~~приделать логгинг~~
1. ~~сделать нормальную обработку ошибок если не получается изменить рекорды~~
2. ~~сделать проверки~~
        ~~интервейс Service~~
            ~~реализации~~
.
3. сделать бизнеслогику
4. сделать конфигурацию
    интерфейс
    реализацию
    
    
--set.normal
--set.fallback
--get.state
--get.config

cfg
----
debug: true/false

route53:
 zoneId: "sdsdfsdf"
 recordeName: "content.cdn.personaly.bid"

cdnHosts:
 - us-01.cdn.personaly.bid
 - us-02.cdn.personaly.bid
 - eu-01.cdn.personaly.bid
 - jp-01.cdn.personaly.bid
 
sslCheck:
  domains:
    - content.cdn.personaly.bid
    - *.cdn.personaly.bid

httpCheck:
  path: "checks/status.txt"
  code: "200"

fallback:
 - action1
 - action2
 - action3
 ...

normal:
 - action1
 - action2
 - action3
 ...
 
