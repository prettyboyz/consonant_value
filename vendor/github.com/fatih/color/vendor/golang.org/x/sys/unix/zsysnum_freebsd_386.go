
// mksysnum_freebsd.pl
// MACHINE GENERATED BY THE ABOVE COMMAND; DO NOT EDIT

// +build 386,freebsd

package unix

const (
	// SYS_NOSYS = 0;  // { int nosys(void); } syscall nosys_args int
	SYS_EXIT                     = 1   // { void sys_exit(int rval); } exit \
	SYS_FORK                     = 2   // { int fork(void); }
	SYS_READ                     = 3   // { ssize_t read(int fd, void *buf, \
	SYS_WRITE                    = 4   // { ssize_t write(int fd, const void *buf, \
	SYS_OPEN                     = 5   // { int open(char *path, int flags, int mode); }
	SYS_CLOSE                    = 6   // { int close(int fd); }
	SYS_WAIT4                    = 7   // { int wait4(int pid, int *status, \
	SYS_LINK                     = 9   // { int link(char *path, char *link); }
	SYS_UNLINK                   = 10  // { int unlink(char *path); }
	SYS_CHDIR                    = 12  // { int chdir(char *path); }
	SYS_FCHDIR                   = 13  // { int fchdir(int fd); }
	SYS_MKNOD                    = 14  // { int mknod(char *path, int mode, int dev); }
	SYS_CHMOD                    = 15  // { int chmod(char *path, int mode); }
	SYS_CHOWN                    = 16  // { int chown(char *path, int uid, int gid); }
	SYS_OBREAK                   = 17  // { int obreak(char *nsize); } break \
	SYS_GETPID                   = 20  // { pid_t getpid(void); }
	SYS_MOUNT                    = 21  // { int mount(char *type, char *path, \
	SYS_UNMOUNT                  = 22  // { int unmount(char *path, int flags); }
	SYS_SETUID                   = 23  // { int setuid(uid_t uid); }
	SYS_GETUID                   = 24  // { uid_t getuid(void); }
	SYS_GETEUID                  = 25  // { uid_t geteuid(void); }
	SYS_PTRACE                   = 26  // { int ptrace(int req, pid_t pid, \
	SYS_RECVMSG                  = 27  // { int recvmsg(int s, struct msghdr *msg, \
	SYS_SENDMSG                  = 28  // { int sendmsg(int s, struct msghdr *msg, \
	SYS_RECVFROM                 = 29  // { int recvfrom(int s, caddr_t buf, \
	SYS_ACCEPT                   = 30  // { int accept(int s, \
	SYS_GETPEERNAME              = 31  // { int getpeername(int fdes, \
	SYS_GETSOCKNAME              = 32  // { int getsockname(int fdes, \
	SYS_ACCESS                   = 33  // { int access(char *path, int amode); }
	SYS_CHFLAGS                  = 34  // { int chflags(const char *path, u_long flags); }
	SYS_FCHFLAGS                 = 35  // { int fchflags(int fd, u_long flags); }
	SYS_SYNC                     = 36  // { int sync(void); }
	SYS_KILL                     = 37  // { int kill(int pid, int signum); }
	SYS_GETPPID                  = 39  // { pid_t getppid(void); }
	SYS_DUP                      = 41  // { int dup(u_int fd); }
	SYS_PIPE                     = 42  // { int pipe(void); }
	SYS_GETEGID                  = 43  // { gid_t getegid(void); }
	SYS_PROFIL                   = 44  // { int profil(caddr_t samples, size_t size, \
	SYS_KTRACE                   = 45  // { int ktrace(const char *fname, int ops, \
	SYS_GETGID                   = 47  // { gid_t getgid(void); }
	SYS_GETLOGIN                 = 49  // { int getlogin(char *namebuf, u_int \
	SYS_SETLOGIN                 = 50  // { int setlogin(char *namebuf); }
	SYS_ACCT                     = 51  // { int acct(char *path); }
	SYS_SIGALTSTACK              = 53  // { int sigaltstack(stack_t *ss, \
	SYS_IOCTL                    = 54  // { int ioctl(int fd, u_long com, \
	SYS_REBOOT                   = 55  // { int reboot(int opt); }
	SYS_REVOKE                   = 56  // { int revoke(char *path); }
	SYS_SYMLINK                  = 57  // { int symlink(char *path, char *link); }
	SYS_READLINK                 = 58  // { ssize_t readlink(char *path, char *buf, \
	SYS_EXECVE                   = 59  // { int execve(char *fname, char **argv, \
	SYS_UMASK                    = 60  // { int umask(int newmask); } umask umask_args \
	SYS_CHROOT                   = 61  // { int chroot(char *path); }
	SYS_MSYNC                    = 65  // { int msync(void *addr, size_t len, \
	SYS_VFORK                    = 66  // { int vfork(void); }
	SYS_SBRK                     = 69  // { int sbrk(int incr); }
	SYS_SSTK                     = 70  // { int sstk(int incr); }
	SYS_OVADVISE                 = 72  // { int ovadvise(int anom); } vadvise \
	SYS_MUNMAP                   = 73  // { int munmap(void *addr, size_t len); }
	SYS_MPROTECT                 = 74  // { int mprotect(const void *addr, size_t len, \
	SYS_MADVISE                  = 75  // { int madvise(void *addr, size_t len, \
	SYS_MINCORE                  = 78  // { int mincore(const void *addr, size_t len, \
	SYS_GETGROUPS                = 79  // { int getgroups(u_int gidsetsize, \
	SYS_SETGROUPS                = 80  // { int setgroups(u_int gidsetsize, \
	SYS_GETPGRP                  = 81  // { int getpgrp(void); }
	SYS_SETPGID                  = 82  // { int setpgid(int pid, int pgid); }
	SYS_SETITIMER                = 83  // { int setitimer(u_int which, struct \
	SYS_SWAPON                   = 85  // { int swapon(char *name); }
	SYS_GETITIMER                = 86  // { int getitimer(u_int which, \
	SYS_GETDTABLESIZE            = 89  // { int getdtablesize(void); }
	SYS_DUP2                     = 90  // { int dup2(u_int from, u_int to); }
	SYS_FCNTL                    = 92  // { int fcntl(int fd, int cmd, long arg); }
	SYS_SELECT                   = 93  // { int select(int nd, fd_set *in, fd_set *ou, \
	SYS_FSYNC                    = 95  // { int fsync(int fd); }
	SYS_SETPRIORITY              = 96  // { int setpriority(int which, int who, \
	SYS_SOCKET                   = 97  // { int socket(int domain, int type, \
	SYS_CONNECT                  = 98  // { int connect(int s, caddr_t name, \
	SYS_GETPRIORITY              = 100 // { int getpriority(int which, int who); }
	SYS_BIND                     = 104 // { int bind(int s, caddr_t name, \
	SYS_SETSOCKOPT               = 105 // { int setsockopt(int s, int level, int name, \
	SYS_LISTEN                   = 106 // { int listen(int s, int backlog); }
	SYS_GETTIMEOFDAY             = 116 // { int gettimeofday(struct timeval *tp, \
	SYS_GETRUSAGE                = 117 // { int getrusage(int who, \
	SYS_GETSOCKOPT               = 118 // { int getsockopt(int s, int level, int name, \
	SYS_READV                    = 120 // { int readv(int fd, struct iovec *iovp, \
	SYS_WRITEV                   = 121 // { int writev(int fd, struct iovec *iovp, \
	SYS_SETTIMEOFDAY             = 122 // { int settimeofday(struct timeval *tv, \
	SYS_FCHOWN                   = 123 // { int fchown(int fd, int uid, int gid); }
	SYS_FCHMOD                   = 124 // { int fchmod(int fd, int mode); }
	SYS_SETREUID                 = 126 // { int setreuid(int ruid, int euid); }
	SYS_SETREGID                 = 127 // { int setregid(int rgid, int egid); }
	SYS_RENAME                   = 128 // { int rename(char *from, char *to); }
	SYS_FLOCK                    = 131 // { int flock(int fd, int how); }
	SYS_MKFIFO                   = 132 // { int mkfifo(char *path, int mode); }
	SYS_SENDTO                   = 133 // { int sendto(int s, caddr_t buf, size_t len, \
	SYS_SHUTDOWN                 = 134 // { int shutdown(int s, int how); }
	SYS_SOCKETPAIR               = 135 // { int socketpair(int domain, int type, \
	SYS_MKDIR                    = 136 // { int mkdir(char *path, int mode); }
	SYS_RMDIR                    = 137 // { int rmdir(char *path); }
	SYS_UTIMES                   = 138 // { int utimes(char *path, \
	SYS_ADJTIME                  = 140 // { int adjtime(struct timeval *delta, \
	SYS_SETSID                   = 147 // { int setsid(void); }
	SYS_QUOTACTL                 = 148 // { int quotactl(char *path, int cmd, int uid, \
	SYS_LGETFH                   = 160 // { int lgetfh(char *fname, \
	SYS_GETFH                    = 161 // { int getfh(char *fname, \
	SYS_SYSARCH                  = 165 // { int sysarch(int op, char *parms); }
	SYS_RTPRIO                   = 166 // { int rtprio(int function, pid_t pid, \
	SYS_FREEBSD6_PREAD           = 173 // { ssize_t freebsd6_pread(int fd, void *buf, \
	SYS_FREEBSD6_PWRITE          = 174 // { ssize_t freebsd6_pwrite(int fd, \
	SYS_SETFIB                   = 175 // { int setfib(int fibnum); }
	SYS_NTP_ADJTIME              = 176 // { int ntp_adjtime(struct timex *tp); }
	SYS_SETGID                   = 181 // { int setgid(gid_t gid); }
	SYS_SETEGID                  = 182 // { int setegid(gid_t egid); }
	SYS_SETEUID                  = 183 // { int seteuid(uid_t euid); }
	SYS_STAT                     = 188 // { int stat(char *path, struct stat *ub); }
	SYS_FSTAT                    = 189 // { int fstat(int fd, struct stat *sb); }
	SYS_LSTAT                    = 190 // { int lstat(char *path, struct stat *ub); }
	SYS_PATHCONF                 = 191 // { int pathconf(char *path, int name); }
	SYS_FPATHCONF                = 192 // { int fpathconf(int fd, int name); }
	SYS_GETRLIMIT                = 194 // { int getrlimit(u_int which, \
	SYS_SETRLIMIT                = 195 // { int setrlimit(u_int which, \
	SYS_GETDIRENTRIES            = 196 // { int getdirentries(int fd, char *buf, \
	SYS_FREEBSD6_MMAP            = 197 // { caddr_t freebsd6_mmap(caddr_t addr, \
	SYS_FREEBSD6_LSEEK           = 199 // { off_t freebsd6_lseek(int fd, int pad, \
	SYS_FREEBSD6_TRUNCATE        = 200 // { int freebsd6_truncate(char *path, int pad, \
	SYS_FREEBSD6_FTRUNCATE       = 201 // { int freebsd6_ftruncate(int fd, int pad, \
	SYS___SYSCTL                 = 202 // { int __sysctl(int *name, u_int namelen, \
	SYS_MLOCK                    = 203 // { int mlock(const void *addr, size_t len); }
	SYS_MUNLOCK                  = 204 // { int munlock(const void *addr, size_t len); }
	SYS_UNDELETE                 = 205 // { int undelete(char *path); }
	SYS_FUTIMES                  = 206 // { int futimes(int fd, struct timeval *tptr); }
	SYS_GETPGID                  = 207 // { int getpgid(pid_t pid); }
	SYS_POLL                     = 209 // { int poll(struct pollfd *fds, u_int nfds, \
	SYS_CLOCK_GETTIME            = 232 // { int clock_gettime(clockid_t clock_id, \
	SYS_CLOCK_SETTIME            = 233 // { int clock_settime( \
	SYS_CLOCK_GETRES             = 234 // { int clock_getres(clockid_t clock_id, \
	SYS_KTIMER_CREATE            = 235 // { int ktimer_create(clockid_t clock_id, \
	SYS_KTIMER_DELETE            = 236 // { int ktimer_delete(int timerid); }
	SYS_KTIMER_SETTIME           = 237 // { int ktimer_settime(int timerid, int flags, \
	SYS_KTIMER_GETTIME           = 238 // { int ktimer_gettime(int timerid, struct \
	SYS_KTIMER_GETOVERRUN        = 239 // { int ktimer_getoverrun(int timerid); }
	SYS_NANOSLEEP                = 240 // { int nanosleep(const struct timespec *rqtp, \
	SYS_FFCLOCK_GETCOUNTER       = 241 // { int ffclock_getcounter(ffcounter *ffcount); }
	SYS_FFCLOCK_SETESTIMATE      = 242 // { int ffclock_setestimate( \
	SYS_FFCLOCK_GETESTIMATE      = 243 // { int ffclock_getestimate( \
	SYS_CLOCK_GETCPUCLOCKID2     = 247 // { int clock_getcpuclockid2(id_t id,\
	SYS_NTP_GETTIME              = 248 // { int ntp_gettime(struct ntptimeval *ntvp); }
	SYS_MINHERIT                 = 250 // { int minherit(void *addr, size_t len, \
	SYS_RFORK                    = 251 // { int rfork(int flags); }
	SYS_OPENBSD_POLL             = 252 // { int openbsd_poll(struct pollfd *fds, \
	SYS_ISSETUGID                = 253 // { int issetugid(void); }
	SYS_LCHOWN                   = 254 // { int lchown(char *path, int uid, int gid); }
	SYS_GETDENTS                 = 272 // { int getdents(int fd, char *buf, \
	SYS_LCHMOD                   = 274 // { int lchmod(char *path, mode_t mode); }
	SYS_LUTIMES                  = 276 // { int lutimes(char *path, \
	SYS_NSTAT                    = 278 // { int nstat(char *path, struct nstat *ub); }
	SYS_NFSTAT                   = 279 // { int nfstat(int fd, struct nstat *sb); }
	SYS_NLSTAT                   = 280 // { int nlstat(char *path, struct nstat *ub); }
	SYS_PREADV                   = 289 // { ssize_t preadv(int fd, struct iovec *iovp, \
	SYS_PWRITEV                  = 290 // { ssize_t pwritev(int fd, struct iovec *iovp, \
	SYS_FHOPEN                   = 298 // { int fhopen(const struct fhandle *u_fhp, \
	SYS_FHSTAT                   = 299 // { int fhstat(const struct fhandle *u_fhp, \
	SYS_MODNEXT                  = 300 // { int modnext(int modid); }
	SYS_MODSTAT                  = 301 // { int modstat(int modid, \
	SYS_MODFNEXT                 = 302 // { int modfnext(int modid); }
	SYS_MODFIND                  = 303 // { int modfind(const char *name); }
	SYS_KLDLOAD                  = 304 // { int kldload(const char *file); }
	SYS_KLDUNLOAD                = 305 // { int kldunload(int fileid); }
	SYS_KLDFIND                  = 306 // { int kldfind(const char *file); }
	SYS_KLDNEXT                  = 307 // { int kldnext(int fileid); }
	SYS_KLDSTAT                  = 308 // { int kldstat(int fileid, struct \
	SYS_KLDFIRSTMOD              = 309 // { int kldfirstmod(int fileid); }
	SYS_GETSID                   = 310 // { int getsid(pid_t pid); }
	SYS_SETRESUID                = 311 // { int setresuid(uid_t ruid, uid_t euid, \
	SYS_SETRESGID                = 312 // { int setresgid(gid_t rgid, gid_t egid, \
	SYS_YIELD                    = 321 // { int yield(void); }
	SYS_MLOCKALL                 = 324 // { int mlockall(int how); }
	SYS_MUNLOCKALL               = 325 // { int munlockall(void); }
	SYS___GETCWD                 = 326 // { int __getcwd(char *buf, u_int buflen); }
	SYS_SCHED_SETPARAM           = 327 // { int sched_setparam (pid_t pid, \
	SYS_SCHED_GETPARAM           = 328 // { int sched_getparam (pid_t pid, struct \
	SYS_SCHED_SETSCHEDULER       = 329 // { int sched_setscheduler (pid_t pid, int \
	SYS_SCHED_GETSCHEDULER       = 330 // { int sched_getscheduler (pid_t pid); }
	SYS_SCHED_YIELD              = 331 // { int sched_yield (void); }
	SYS_SCHED_GET_PRIORITY_MAX   = 332 // { int sched_get_priority_max (int policy); }
	SYS_SCHED_GET_PRIORITY_MIN   = 333 // { int sched_get_priority_min (int policy); }
	SYS_SCHED_RR_GET_INTERVAL    = 334 // { int sched_rr_get_interval (pid_t pid, \
	SYS_UTRACE                   = 335 // { int utrace(const void *addr, size_t len); }
	SYS_KLDSYM                   = 337 // { int kldsym(int fileid, int cmd, \
	SYS_JAIL                     = 338 // { int jail(struct jail *jail); }
	SYS_SIGPROCMASK              = 340 // { int sigprocmask(int how, \
	SYS_SIGSUSPEND               = 341 // { int sigsuspend(const sigset_t *sigmask); }
	SYS_SIGPENDING               = 343 // { int sigpending(sigset_t *set); }
	SYS_SIGTIMEDWAIT             = 345 // { int sigtimedwait(const sigset_t *set, \
	SYS_SIGWAITINFO              = 346 // { int sigwaitinfo(const sigset_t *set, \
	SYS___ACL_GET_FILE           = 347 // { int __acl_get_file(const char *path, \
	SYS___ACL_SET_FILE           = 348 // { int __acl_set_file(const char *path, \
	SYS___ACL_GET_FD             = 349 // { int __acl_get_fd(int filedes, \
	SYS___ACL_SET_FD             = 350 // { int __acl_set_fd(int filedes, \
	SYS___ACL_DELETE_FILE        = 351 // { int __acl_delete_file(const char *path, \
	SYS___ACL_DELETE_FD          = 352 // { int __acl_delete_fd(int filedes, \