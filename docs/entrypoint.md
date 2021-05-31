# 运用gosu
*   创建group和普通账号，不要使用root账号启动进程
*   如果普通账号权限不够用，建议使用gosu来提升权限，而不是sudo
*   entrypoint.sh脚本在执行的时候也是个进程，启动业务进程的时候，在命令前面加上exec，这样新的进程就会取代entrypoint.sh的进程，得到1号PID
*   exec "$@"是个保底的逻辑，如果entrypoint.sh的入参在整个脚本中都没有被执行，那么exec "$@"会把入参执行一遍，如果前面执行过了，这一行就不起作用，这个命令的细节在Stack Overflow上有详细的描述，如下图，地址是：https://stackoverflow.com/questions/39082768/what-does-set-e-and-exec-do-for-docker-entrypoint-scripts
# gosu与sudo
*   gosu启动命令时只有一个进程，所以docker容器启动时使用gosu，那么该进程可以做到PID等于1
*   sudo启动命令时先创建sudo进程，然后该进程作为父进程去创建子进程，1号PID被sudo进程占据