// mkerrors.sh
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

// +build mips,linux

// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs -- _const.go

package unix

import "syscall"

const (
	AF_ALG                           = 0x26
	AF_APPLETALK                     = 0x5
	AF_ASH                           = 0x12
	AF_ATMPVC                        = 0x8
	AF_ATMSVC                        = 0x14
	AF_AX25                          = 0x3
	AF_BLUETOOTH                     = 0x1f
	AF_BRIDGE                        = 0x7
	AF_CAIF                          = 0x25
	AF_CAN                           = 0x1d
	AF_DECnet                        = 0xc
	AF_ECONET                        = 0x13
	AF_FILE                          = 0x1
	AF_IEEE802154                    = 0x24
	AF_INET                          = 0x2
	AF_INET6                         = 0xa
	AF_IPX                           = 0x4
	AF_IRDA                          = 0x17
	AF_ISDN                          = 0x22
	AF_IUCV                          = 0x20
	AF_KEY                           = 0xf
	AF_LLC                           = 0x1a
	AF_LOCAL                         = 0x1
	AF_MAX                           = 0x27
	AF_NETBEUI                       = 0xd
	AF_NETLINK                       = 0x10
	AF_NETROM                        = 0x6
	AF_PACKET                        = 0x11
	AF_PHONET                        = 0x23
	AF_PPPOX                         = 0x18
	AF_RDS                           = 0x15
	AF_ROSE                          = 0xb
	AF_ROUTE                         = 0x10
	AF_RXRPC                         = 0x21
	AF_SECURITY                      = 0xe
	AF_SNA                           = 0x16
	AF_TIPC                          = 0x1e
	AF_UNIX                          = 0x1
	AF_UNSPEC                        = 0x0
	AF_WANPIPE                       = 0x19
	AF_X25                           = 0x9
	ALG_OP_DECRYPT                   = 0x0
	ALG_OP_ENCRYPT                   = 0x1
	ALG_SET_AEAD_ASSOCLEN            = 0x4
	ALG_SET_AEAD_AUTHSIZE            = 0x5
	ALG_SET_IV                       = 0x2
	ALG_SET_KEY                      = 0x1
	ALG_SET_OP                       = 0x3
	ARPHRD_ADAPT                     = 0x108
	ARPHRD_APPLETLK                  = 0x8
	ARPHRD_ARCNET                    = 0x7
	ARPHRD_ASH                       = 0x30d
	ARPHRD_ATM                       = 0x13
	ARPHRD_AX25                      = 0x3
	ARPHRD_BIF                       = 0x307
	ARPHRD_CAIF                      = 0x336
	ARPHRD_CAN                       = 0x118
	ARPHRD_CHAOS                     = 0x5
	ARPHRD_CISCO                     = 0x201
	ARPHRD_CSLIP                     = 0x101
	ARPHRD_CSLIP6                    = 0x103
	ARPHRD_DDCMP                     = 0x205
	ARPHRD_DLCI                      = 0xf
	ARPHRD_ECONET                    = 0x30e
	ARPHRD_EETHER                    = 0x2
	ARPHRD_ETHER                     = 0x1
	ARPHRD_EUI64                     = 0x1b
	ARPHRD_FCAL                      = 0x311
	ARPHRD_FCFABRIC                  = 0x313
	ARPHRD_FCPL                      = 0x312
	ARPHRD_FCPP                      = 0x310
	ARPHRD_FDDI                      = 0x306
	ARPHRD_FRAD                      = 0x302
	ARPHRD_HDLC                      = 0x201
	ARPHRD_HIPPI                     = 0x30c
	ARPHRD_HWX25                     = 0x110
	ARPHRD_IEEE1394                  = 0x18
	ARPHRD_IEEE802                   = 0x6
	ARPHRD_IEEE80211                 = 0x321
	ARPHRD_IEEE80211_PRISM           = 0x322
	ARPHRD_IEEE80211_RADIOTAP        = 0x323
	ARPHRD_IEEE802154                = 0x324
	ARPHRD_IEEE802_TR                = 0x320
	ARPHRD_INFINIBAND                = 0x20
	ARPHRD_IPDDP                     = 0x309
	ARPHRD_IPGRE                     = 0x30a
	ARPHRD_IRDA                      = 0x30f
	ARPHRD_LAPB                      = 0x204
	ARPHRD_LOCALTLK                  = 0x305
	ARPHRD_LOOPBACK                  = 0x304
	ARPHRD_METRICOM                  = 0x17
	ARPHRD_NETROM                    = 0x0
	ARPHRD_NONE                      = 0xfffe
	ARPHRD_PHONET                    = 0x334
	ARPHRD_PHONET_PIPE               = 0x335
	ARPHRD_PIMREG                    = 0x30b
	ARPHRD_PPP                       = 0x200
	ARPHRD_PRONET                    = 0x4
	ARPHRD_RAWHDLC                   = 0x206
	ARPHRD_ROSE                      = 0x10e
	ARPHRD_RSRVD                     = 0x104
	ARPHRD_SIT                       = 0x308
	ARPHRD_SKIP                      = 0x303
	ARPHRD_SLIP                      = 0x100
	ARPHRD_SLIP6                     = 0x102
	ARPHRD_TUNNEL                    = 0x300
	ARPHRD_TUNNEL6                   = 0x301
	ARPHRD_VOID                      = 0xffff
	ARPHRD_X25                       = 0x10f
	B0                               = 0x0
	B1000000                         = 0x1008
	B110                             = 0x3
	B115200                          = 0x1002
	B1152000                         = 0x1009
	B1200                            = 0x9
	B134                             = 0x4
	B150                             = 0x5
	B1500000                         = 0x100a
	B1800                            = 0xa
	B19200                           = 0xe
	B200                             = 0x6
	B2000000                         = 0x100b
	B230400                          = 0x1003
	B2400                            = 0xb
	B2500000                         = 0x100c
	B300                             = 0x7
	B3000000                         = 0x100d
	B3500000                         = 0x100e
	B38400                           = 0xf
	B4000000                         = 0x100f
	B460800                          = 0x1004
	B4800                            = 0xc
	B50                              = 0x1
	B500000                          = 0x1005
	B57600                           = 0x1001
	B576000                          = 0x1006
	B600                             = 0x8
	B75                              = 0x2
	B921600                          = 0x1007
	B9600                            = 0xd
	BOTHER                           = 0x1000
	BPF_A                            = 0x10
	BPF_ABS                          = 0x20
	BPF_ADD                          = 0x0
	BPF_ALU                          = 0x4
	BPF_AND                          = 0x50
	BPF_B                            = 0x10
	BPF_DIV                          = 0x30
	BPF_H                            = 0x8
	BPF_IMM                          = 0x0
	BPF_IND                          = 0x40
	BPF_JA                           = 0x0
	BPF_JEQ                          = 0x10
	BPF_JGE                          = 0x30
	BPF_JGT                          = 0x20
	BPF_JMP                          = 0x5
	BPF_JSET                         = 0x40
	BPF_K                            = 0x0
	BPF_LD                           = 0x0
	BPF_LDX                          = 0x1
	BPF_LEN                          = 0x80
	BPF_LSH                          = 0x60
	BPF_MAJOR_VERSION                = 0x1
	BPF_MAXINSNS                     = 0x1000
	BPF_MEM                          = 0x60
	BPF_MEMWORDS                     = 0x10
	BPF_MINOR_VERSION                = 0x1
	BPF_MISC                         = 0x7
	BPF_MSH                          = 0xa0
	BPF_MUL                          = 0x20
	BPF_NEG                          = 0x80
	BPF_OR                           = 0x40
	BPF_RET                          = 0x6
	BPF_RSH                          = 0x70
	BPF_ST                           = 0x2
	BPF_STX                          = 0x3
	BPF_SUB                          = 0x10
	BPF_TAX                          = 0x0
	BPF_TXA                          = 0x80
	BPF_W                            = 0x0
	BPF_X                            = 0x8
	BRKINT                           = 0x2
	BS0                              = 0x0
	BS1                              = 0x2000
	BSDLY                            = 0x2000
	CAN_BCM                          = 0x2
	CAN_EFF_FLAG                     = 0x80000000
	CAN_EFF_MASK                     = 0x1fffffff
	CAN_ERR_FLAG                     =