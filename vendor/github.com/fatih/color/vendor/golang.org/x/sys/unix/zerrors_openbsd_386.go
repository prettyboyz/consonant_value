// mkerrors.sh -m32
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

// +build 386,openbsd

// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs -- -m32 _const.go

package unix

import "syscall"

const (
	AF_APPLETALK                      = 0x10
	AF_BLUETOOTH                      = 0x20
	AF_CCITT                          = 0xa
	AF_CHAOS                          = 0x5
	AF_CNT                            = 0x15
	AF_COIP                           = 0x14
	AF_DATAKIT                        = 0x9
	AF_DECnet                         = 0xc
	AF_DLI                            = 0xd
	AF_E164                           = 0x1a
	AF_ECMA                           = 0x8
	AF_ENCAP                          = 0x1c
	AF_HYLINK                         = 0xf
	AF_IMPLINK                        = 0x3
	AF_INET                           = 0x2
	AF_INET6                          = 0x18
	AF_IPX                            = 0x17
	AF_ISDN                           = 0x1a
	AF_ISO                            = 0x7
	AF_KEY                            = 0x1e
	AF_LAT                            = 0xe
	AF_LINK                           = 0x12
	AF_LOCAL                          = 0x1
	AF_MAX                            = 0x24
	AF_MPLS                           = 0x21
	AF_NATM                           = 0x1b
	AF_NS                             = 0x6
	AF_OSI                            = 0x7
	AF_PUP                            = 0x4
	AF_ROUTE                          = 0x11
	AF_SIP                            = 0x1d
	AF_SNA                            = 0xb
	AF_UNIX                           = 0x1
	AF_UNSPEC                         = 0x0
	ARPHRD_ETHER                      = 0x1
	ARPHRD_FRELAY                     = 0xf
	ARPHRD_IEEE1394                   = 0x18
	ARPHRD_IEEE802                    = 0x6
	B0                                = 0x0
	B110                              = 0x6e
	B115200                           = 0x1c200
	B1200                             = 0x4b0
	B134                              = 0x86
	B14400                            = 0x3840
	B150                              = 0x96
	B1800                             = 0x708
	B19200                            = 0x4b00
	B200                              = 0xc8
	B230400                           = 0x38400
	B2400                             = 0x960
	B28800                            = 0x7080
	B300                              = 0x12c
	B38400                            = 0x9600
	B4800                             = 0x12c0
	B50                               = 0x32
	B57600                            = 0xe100
	B600                              = 0x258
	B7200                             = 0x1c20
	B75                               = 0x4b
	B76800                            = 0x12c00
	B9600                             = 0x2580
	BIOCFLUSH                         = 0x20004268
	BIOCGBLEN                         = 0x40044266
	BIOCGDIRFILT                      = 0x4004427c
	BIOCGDLT                          = 0x4004426a
	BIOCGDLTLIST                      = 0xc008427b
	BIOCGETIF                         = 0x4020426b
	BIOCGFILDROP                      = 0x40044278
	BIOCGHDRCMPLT                     = 0x40044274
	BIOCGRSIG                         = 0x40044273
	BIOCGRTIMEOUT                     = 0x400c426e
	BIOCGSTATS                        = 0x4008426f
	BIOCIMMEDIATE                     = 0x80044270
	BIOCLOCK                          = 0x20004276
	BIOCPROMISC                       = 0x20004269
	BIOCSBLEN                         = 0xc0044266
	BIOCSDIRFILT                      = 0x8004427d
	BIOCSDLT                          = 0x8004427a
	BIOCSETF                          = 0x80084267
	BIOCSETIF                         = 0x8020426c
	BIOCSETWF                         = 0x80084277
	BIOCSFILDROP                      = 0x80044279
	BIOCSHDRCMPLT                     = 0x80044275
	BIOCSRSIG                         = 0x80044272
	BIOCSRTIMEOUT                     = 0x800c426d
	BIOCVERSION                       = 0x40044271
	BPF_A                             = 0x10
	BPF_ABS                           = 0x20
	BPF_ADD                           = 0x0
	BPF_ALIGNMENT                     = 0x4
	BPF_ALU                           = 0x4
	BPF_AND                           = 0x50
	BPF_B                             = 0x10
	BPF_DIRECTION_IN                  = 0x1
	BPF_DIRECTION_OUT                 = 0x2
	BPF_DIV                           = 0x30
	BPF_H                             = 0x8
	BPF_IMM                           = 0x0
	BPF_IND                           = 0x40
	BPF_JA                            = 0x0
	BPF_JEQ                           = 0x10
	BPF_JGE                           = 0x30
	BPF_JGT                           = 0x20
	BPF_JMP                           = 0x5
	BPF_JSET                          = 0x40
	BPF_K                             = 0x0
	BPF_LD                            = 0x0
	BPF_LDX                           = 0x1
	BPF_LEN                           = 0x80
	BPF_LSH                           = 0x60
	BPF_MAJOR_VERSION                 = 0x1
	BPF_MAXBUFSIZE                    = 0x200000
	BPF_MAXINSNS                      = 0x200
	BPF_MEM                           = 0x60
	BPF_MEMWORDS                      = 0x10
	BPF_MINBUFSIZE                    = 0x20
	BPF_MINOR_VERSION                 = 0x1
	BPF_MISC                          = 0x7
	BPF_MSH                           = 0xa0
	BPF_MUL                           = 0x20
	BPF_NEG                           = 0x80
	BPF_OR                            = 0x40
	BPF_RELEASE                       = 0x30bb6
	BPF_RET                           = 0x6
	BPF_RSH                           = 0x70
	BPF_ST                            = 0x2
	BPF_STX                           = 0x3
	BPF_SUB                           = 0x10
	BPF_TAX                           = 0x0
	BPF_TXA                           = 0x80
	BPF_W                             = 0x0
	BPF_X                             = 0x8
	BRKINT                            = 0x2
	CFLUSH                            = 0xf
	CLOCAL                            = 0x8000
	CREAD                             = 0x800
	CS5                               = 0x0
	CS6                               = 0x100
	CS7                               = 0x200
	CS8                               = 0x300
	CSIZE                             = 0x300
	CSTART                            = 0x11
	CSTATUS                           = 0xff
	CSTOP                             = 0x13
	CSTOPB                            = 0x400
	CSUSP                             = 0x1a
	CTL_MAXNAME                       = 0xc
	CTL_NET                           = 0x4
	DIOCOSFPFLUSH                     = 0x2000444e
	DLT_ARCNET                        = 0x7
	DLT_ATM_RFC1483                   = 0xb
	DLT_AX25                          = 0x3
	DLT_CHAOS                         = 0x5
	DLT_C_HDLC                        = 0x68
	DLT_EN10MB                        = 0x1
	DLT_EN3MB                         = 0x2
	DLT_ENC                           = 0xd
	DLT_FDDI                          = 0xa
	DLT_IEEE802                       = 0x6
	DLT_IEEE802_11                    = 0x69
	DLT_IEEE802_11_RADIO              = 0x7f
	DLT_LOOP                          = 0xc
	DLT_MPLS                          = 0xdb
	DLT_NULL                          = 0x0
	DLT_PFLOG                         = 0x75
	DLT_PFSYNC                        = 0x12
	DLT_PPP                           = 0x9
	DLT_PPP_BSDOS                     = 0x10
	DLT_PPP_ETHER                     = 0x33
	DLT_PPP_SERIAL                    = 0x32
	DLT_PRONET                        = 0x4
	DLT_RAW                           = 0xe
	DLT_SLIP                          = 0x8
	DLT_SLIP_BSDOS                    = 0xf
	DT_BLK                            = 0x6
	DT_CHR                            = 0x2
	DT_DIR                            = 0x4
	DT_FIFO                           = 0x1
	DT_LNK                            = 0xa
	DT_REG                            = 0x8
	DT_SOCK                           = 0xc
	DT_UNKNOWN                        = 0x0
	ECHO                              = 0x8
	ECHOCTL                           = 0x40
	ECHOE                             = 0x2
	ECHOK                             = 0x4
	ECHOKE                            = 0x1
	ECHONL                            = 0x10
	ECHOPRT                           = 0x20
	EMT_TAGOVF                        = 0x1
	EMUL_ENABLED                      = 0x1
	EMUL_NATIVE                       = 0x2
	ENDRUNDISC                        = 0x9
	ETHERMIN                          = 0x2e
	ETHERMTU                          = 0x5dc
	ETHERTYPE_8023                    = 0x4
	ETHERTYPE_AARP                    = 0x80f3
	ETHERTYPE_ACCTON                  = 0x8390
	ETHERTYPE_AEONIC                  = 0x8036
	ETHERTYPE_ALPHA                   = 0x814a
	ETHERTYPE_AMBER                   = 0x6008
	ETHERTYPE_AMOEBA                  = 0x8145
	ETHERTYPE_AOE                     = 0x88a2
	ETHERTYPE_APOLLO                  = 0x80f7
	ETHERTYPE_APOLLODOMAIN            = 0x8019
	ETHERTYPE_APPLETALK               = 0x809b
	ETHERTYPE_APPLITEK                = 0x80c7
	ETHERTYPE_ARGONAUT                = 0x803a
	ETHERTYPE_ARP                     = 0x806
	ETHERTYPE_AT                      = 0x809b
	ETHERTYPE_ATALK                   = 0x809b
	ETHERTYPE_ATOMIC                  = 0x86df
	ETHERTYPE_ATT                     = 0x8069
	ETHERTYPE_ATTSTANFORD             = 0x8008
	ETHERTYPE_AUTOPHON                = 0x806a
	ETHERTYPE_AXIS                    = 0x8856
	ETHERTYPE_BCLOOP                  = 0x9003
	ETHERTYPE_BOFL                    = 0x8102
	ETHERTYPE_CABLETRON               = 0x7034
	ETHERTYPE_CHAOS                   = 0x804
	ETHERTYPE_COMDESIGN               = 0x806c
	ETHERTYPE_COMPUGRAPHIC            = 0x806d
	ETHERTYPE_COUNTERPOINT            = 0x8062
	ETHERTYPE_CRONUS                  = 0x8004
	ETHERTYPE_CRONUSVLN               = 0x8003
	ETHERTYPE_DCA                     = 0x1234
	ETHERTYPE_DDE                     = 0x807b
	ETHERTYPE_DEBNI                   = 0xaaaa
	ETHERTYPE_DECAM                   = 0x8048
	ETHERTYPE_DECCUST                 = 0x6006
	ETHERTYPE_DECDIAG                 = 0x6005
	ETHERTYPE_DECDNS                  = 0x803c
	ETHERTYPE_DECDTS                  = 0x803e
	ETHERTYPE_DECEXPER                = 0x6000
	ETHERTYPE_DECLAST                 = 0x8041
	ETHERTYPE_DECLTM                  = 0x803f
	ETHERTYPE_DECMUMPS                = 0x6009
	ETHERTYPE_DECNETBIOS              = 0x8040
	ETHERTYPE_DELTACON                = 0x86de
	ETHERTYPE_DIDDLE                  = 0x4321
	ETHERTYPE_DLOG1                   = 0x660
	ETHERTYPE_DLOG2                   = 0x661
	ETHERTYPE_DN                      = 0x6003
	ETHERTYPE_DOGFIGHT                = 0x1989
	ETHERTYPE_DSMD                    = 0x8039
	ETHERTYPE_ECMA                    = 0x803
	ETHERTYPE_ENCRYPT                 = 0x803d
	ETHERTYPE_ES                      = 0x805d
	ETHERTYPE_EXCELAN                 = 0x8010
	ETHERTYPE_EXPERDATA               = 0x8049
	ETHERTYPE_FLIP                    = 0x8146
	ETHERTYPE_FLOWCONTROL             = 0x8808
	ETHERTYPE_FRARP                   = 0x808
	ETHERTYPE_GENDYN                  = 0x8068
	ETHERTYPE_HAYES                   = 0x8130
	ETHERTYPE_HIPPI_FP                = 0x8180
	ETHERTYPE_HITACHI                 = 0x8820
	ETHERTYPE_HP                      = 0x8005
	ETHERTYPE_IEEEPUP                 = 0xa00
	ETHERTYPE_IEEEPUPAT               = 0xa01
	ETHERTYPE_IMLBL                   = 0x4c42
	ETHERTYPE_IMLBLDIAG               = 0x424c
	ETHERTYPE_IP                      = 0x800
	ETHERTYPE_IPAS                    =