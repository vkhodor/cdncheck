3. сделать отправку сообщения в slack/telegram - если происходит переключение должно отправляться сообщение в телегу, слэк, почту...
4. сделать конфигурацию
    интерфейс
    реализацию
    
----------------------------------    
--set.normal
--set.fallback
--get.state
----------------------------------
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
 
