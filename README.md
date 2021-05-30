# github.com/crochee/proxy-go

The goal is to build a high-performance general-purpose gateway

## 目录结构
├─asset  
├─build  
│  ├─ci  
│  └─package  
│      └─proxy  
│          └─cert  
├─cmd  
│  ├─proxy  
│  └─web  
├─conf  
├─config  
│  └─dynamic  
├─deployment  
├─docs  
├─internal  
├─log  
├─pkg  
│  ├─filecontent  
│  ├─logger  
│  ├─metrics  
│  ├─middleware  
│  │  ├─accesslog  
│  │  ├─balance  
│  │  ├─circuitbreaker  
│  │  ├─cros  
│  │  ├─metric  
│  │  ├─ratelimit  
│  │  ├─recovery  
│  │  ├─retry  
│  │  └─trace  
│  ├─proxy  
│  │  ├─httpx  
│  │  └─tcpx  
│  ├─router  
│  ├─routine  
│  ├─selector  
│  ├─tlsx  
│  ├─tracex  
│  │  └─jaeger  
│  ├─transport  
│  │  ├─httpx  
│  │  ├─pprofx  
│  │  └─prometheusx  
│  └─writer  
├─script  
├─test  
│  └─data  
├─version  
└─website  
