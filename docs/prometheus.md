# 普罗米修斯
参考https://prometheus.io/docs/prometheus/latest/installation/
## 指标类型
### Counter 计数器
counter 是一个累计的指标，代表一个单调递增的计数器，它的值只会增加或在重启时重置为零。例如，你可以使用 counter 来代表服务过的请求数，完成的任务数，或者错误的次数。
### Gauge 计量器
gauge 是代表一个数值类型的指标，它的值可以增或减。gauge 通常用于一些度量的值例如温度或是当前内存使用，也可以用于一些可以增减的“计数”，如正在运行的 Goroutine 个数。
### histogram 分布图
histogram 对观测值（类似请求延迟或回复包大小）进行采样，并用一些可配置的桶来计数。它也会给出一个所有观测值的总和。
### Summary 摘要
跟 histogram 类似，summary 也对观测值（类似请求延迟或回复包大小）进行采样。同时它会给出一个总数以及所有观测值的总和，它在一个滑动的时间窗口上计算可配置的分位数。
