// mkerrors.sh -m64
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

// +build amd64,darwin

// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs -- -m64 _const.go

package unix

import "syscall"

const (
	AF_APPLETALK                      = 0x10
	AF_CCITT                          = 0xa
	AF_CHAOS                          = 0x5
	AF_CNT                            = 0x15
	AF_COIP                           = 0x14
	AF_DATAKIT                        = 0x9
	AF_DECnet                         = 0xc
	AF_DLI                            = 0xd
	AF_E164                           = 0x1c
	AF_ECMA                           = 0x8
	AF_HYLINK                         = 0xf
	AF_IEEE80211                      = 0x25
	AF_IMPLINK                        = 0x3
	AF_INET                           = 0x2
	AF_INET6                          = 0x1e
	AF_IPX                            = 0x17
	AF_ISDN                           = 0x1c
	AF_ISO                            = 0x7
	AF_LAT                            = 0xe
	AF_LINK                           = 0x12
	AF_LOCAL                          = 0x1
	AF_MAX                            = 0x28
	AF_NATM                           = 0x1f
	AF_NDRV                           = 0x1b
	AF_NETBIOS                        = 0x21
	AF_NS                             = 0x6
	AF_OSI                            = 0x7
	AF_PPP                            = 0x22
	AF_PUP                            = 0x4
	AF_RESERVED_36                    = 0x24
	AF_ROUTE                          = 0x11
	AF_SIP                            = 0x18
	AF_SNA                            = 0xb
	AF_SYSTEM                         = 0x20
	AF_UNIX                           = 0x1
	AF_UNSPEC                         = 0x0
	AF_UTUN                           = 0x26
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
	BIOCGDLT                          = 0x4004426a
	BIOCGDLTLIST                      = 0xc00c4279
	BIOCGETIF                         = 0x4020426b
	BIOCGHDRCMPLT                     = 0x40044274
	BIOCGRSIG                         = 0x40044272
	BIOCGRTIMEOUT                     = 0x4010426e
	BIOCGSEESENT                      = 0x40044276
	BIOCGSTATS                        = 0x4008426f
	BIOCIMMEDIATE                     = 0x80044270
	BIOCPROMISC                       = 0x20004269
	BIOCSBLEN                         = 0xc0044266
	BIOCSDLT                          = 0x80044278
	BIOCSETF                          = 0x80104267
	BIOCSETFNR                        = 0x8010427e
	BIOCSETIF                         = 0x8020426c
	BIOCSHDRCMPLT                     = 0x80044275
	BIOCSRSIG                         = 0x80044273
	BIOCSRTIMEOUT                     = 0x8010426d
	BIOCSSEESENT                      = 0x80044277
	BIOCVERSION                       = 0x40044271
	BPF_A                             = 0x10
	BPF_ABS                           = 0x20
	BPF_ADD                           = 0x0
	BPF_ALIGNMENT                     = 0x4
	BPF_ALU                           = 0x4
	BPF_AND                           = 0x50
	BPF_B                             = 0x10
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
	BPF_MAXBUFSIZE                    = 0x80000
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
	CSTATUS                           = 0x14
	CSTOP                             = 0x13
	CSTOPB                            = 0x400
	CSUSP                             = 0x1a
	CTL_MAXNAME                       = 0xc
	CTL_NET                           = 0x4
	DLT_A429                          = 0xb8
	DLT_A653_ICM                      = 0xb9
	DLT_AIRONET_HEADER                = 0x78
	DLT_AOS                           = 0xde
	DLT_APPLE_IP_OVER_IEEE1394        = 0x8a
	DLT_ARCNET                        = 0x7
	DLT_ARCNET_LINUX                  = 0x81
	DLT_ATM_CLIP                      = 0x13
	DLT_ATM_RFC1483                   = 0xb
	DLT_AURORA                        = 0x7e
	DLT_AX25                          = 0x3
	DLT_AX25_KISS                     = 0xca
	DLT_BACNET_MS_TP                  = 0xa5
	DLT_BLUETOOTH_HCI_H4              = 0xbb
	DLT_BLUETOOTH_HCI_H4_WITH_PHDR    = 0xc9
	DLT_CAN20B                        = 0xbe
	DLT_CAN_SOCKETCAN                 = 0xe3
	DLT_CHAOS                         = 0x5
	DLT_CHDLC                         = 0x68
	DLT_CISCO_IOS                     = 0x76
	DLT_C_HDLC                        = 0x68
	DLT_C_HDLC_WITH_DIR               = 0xcd
	DLT_DBUS                          = 0xe7
	DLT_DECT                          = 0xdd
	DLT_DOCSIS                        = 0x8f
	DLT_DVB_CI                        = 0xeb
	DLT_ECONET                        = 0x73
	DLT_EN10MB                        = 0x1
	DLT_EN3MB                         = 0x2
	DLT_ENC                           = 0x6d
	DLT_ERF                           = 0xc5
	DLT_ERF_ETH                       = 0xaf
	DLT_ERF_POS                       = 0xb0
	DLT_FC_2                          = 0xe0
	DLT_FC_2_WITH_FRAME_DELIMS        = 0xe1
	DLT_FDDI                          = 0xa
	DLT_FLEXRAY                       = 0xd2
	DLT_FRELAY                        = 0x6b
	DLT_FRELAY_WITH_DIR               = 0xce
	DLT_GCOM_SERIAL                   = 0xad
	DLT_GCOM_T1E1                     = 0xac
	DLT_GPF_F                         = 0xab
	DLT_GPF_T                         = 0xaa
	DLT_GPRS_LLC                      = 0xa9
	DLT_GSMTAP_ABIS                   = 0xda
	DLT_GSMTAP_UM                     = 0xd9
	DLT_HHDLC                         = 0x79
	DLT_IBM_SN                        = 0x92
	DLT_IBM_SP                        = 0x91
	DLT_IEEE802                       = 0x6
	DLT_IEEE802_11                    = 0x69
	DLT_IEEE802_11_RADIO              = 0x7f
	DLT_IEEE802_11_RADIO_AVS          = 0xa3
	DLT_IEEE802_15_4                  = 0xc3
	DLT_IEEE802_15_4_LINUX            = 0xbf
	DLT_IEEE802_15_4_NOFCS            = 0xe6
	DLT_IEEE802_15_4_NONASK_PHY       = 0xd7
	DLT_IEEE802_16_MAC_CPS            = 0xbc
	DLT_IEEE802_16_MAC_CPS_RADIO      = 0xc1
	DLT_IPFILTER                      = 0x74
	DLT_IPMB                          = 0xc7
	DLT_IPMB_LINUX                    = 0xd1
	DLT_IPNET                         = 0xe2
	DLT_IPOIB                         = 0xf2
	DLT_IPV4                          = 0xe4
	DLT_IPV6                          = 0xe5
	DLT_IP_OVER_FC                    = 0x7a
	DLT_JUNIPER_ATM1                  = 0x89
	DLT_JUNIPER_ATM2                  = 0x87
	DLT_JUNIPER_ATM_CEMIC             = 0xee
	DLT_JUNIPER_CHDLC                 = 0xb5
	DLT_JUNIPER_ES                    = 0x84
	DLT_JUNIPER_ETHER                 = 0xb2
	DLT_JUNIPER_FIBRECHANNEL          = 0xea
	DLT_JUNIPER_FRELAY                = 0xb4
	DLT_JUNIPER_GGSN                  = 0x85
	DLT_JUNIPER_ISM                   = 0xc2
	DLT_JUNIPER_MFR                   = 0x86
	DLT_JUNIPER_MLFR                  = 0x83
	DLT_JUNIPER_MLPPP                 = 0x82
	DLT_JUNIPER_MONITOR               = 0xa4
	DLT_JUNIPER_PIC_PEER              = 0xae
	DLT_JUNIPER_PPP                   = 0xb3
	DLT_JUNIPER_PPPOE                 = 0xa7
	DLT_JUNIPER_PPPOE_ATM             = 0xa8
	DLT_JUNIPER_SERVICES              = 0x88
	DLT_JUNIPER_SRX_E2E               = 0xe9
	DLT_JUNIPER_ST                    = 0xc8
	DLT_JUNIPER_VP                    = 0xb7
	DLT_JUNIPER_VS                    = 0xe8
	DLT_LAPB_WITH_DIR                 = 0xcf
	DLT_LAPD                          = 0xcb
	DLT_LIN                           = 0xd4
	DLT_LINUX_EVDEV                   = 0xd8
	DLT_LINUX_IRDA                    = 0x90
	DLT_LINUX_LAPD                    = 0xb1
	DLT_LINUX_PPP_WITHDIRECTION       = 0xa6
	DLT_LINUX_SLL                     = 0x71
	DLT_LOOP                          = 0x6c
	DLT_LTALK                         = 0x72
	DLT_MATCHING_MAX                  = 0xf5
	DLT_MATCHING_MIN                  = 0x68
	DLT_MFR                           = 0xb6
	DLT_MOST                          = 0xd3
	DLT_MPEG_2_TS                     = 0xf3
	DLT_MPLS                          = 0xdb
	DLT_MTP2                          = 0x8c
	DLT_MTP2_WITH_PHDR                = 0x8b
	DLT_MTP3                          = 0x8d
	DLT_MUX27010                      = 0xec
	DLT_NETANALYZER                   = 0xf0
	DLT_NETANALYZER_TRANSPARENT       = 0xf1
	DLT_NFC_LLCP                      = 0xf5
	DLT_NFLOG                         = 0xef
	DLT_NG40                          = 0xf4
	DLT_NULL                          = 0x0
	DLT_PCI_EXP                       = 0x7d
	DLT_PFLOG                         = 0x75
	DLT_PFSYNC                        = 0x12
	DLT_PPI                           = 0xc0
	DLT_PPP                           = 0x9
	DLT_PPP_BSDOS                     = 0x10
	DLT_PPP_ETHER                     = 0x33
	DLT_PPP_PPPD                      = 0xa6
	DLT_PPP_SERIAL                    = 0x32
	DLT_PPP_WITH_DIR                  = 0xcc
	DLT_PPP_WITH_DIRECTION            = 0xa6
	DLT_PRISM_HEADER                  = 0x77
	DLT_PRONET                        = 0x4
	DLT_RAIF1                         = 0xc6
	DLT_RAW                           = 0xc
	DLT_RIO                           = 0x7c
	DLT_SCCP                          = 0x8e
	DLT_SITA                          = 0xc4
	DLT_SLIP                          = 0x8
	DLT_SLIP_BSDOS                    = 0xf
	DLT_STANAG_5066_D_PDU             = 0xed
	DLT_SUNATM                        = 0x7b
	DLT_SYMANTEC_FIREWALL             = 0x63
	DLT_TZSP                          = 0x80
	DLT_USB                           = 0xba
	DLT_USB_LINUX                     = 0xbd
	DLT_USB_LINUX_MMAPPED             = 0xdc
	DLT_USER0                         = 0x93
	DLT_USER1                         = 0x94
	DLT_USER10                        = 0x9d
	DLT_USER11                        = 0x9e
	DLT_USER12                        = 0x9f
	DLT_USER13                        = 0xa0
	DLT_USER14                        = 0xa1
	DLT_USER15                        = 0xa2
	DLT_USER2                         = 0x95
	DLT_USER3                         = 0x96
	DLT_USER4                         = 0x97
	DLT_USER5                         = 0x98
	DLT_USER6                         = 0x99
	DLT_USER7                         = 0x9a
	DLT_USER8                         = 0x9b
	DLT_USER9                         = 0x9c
	DLT_WIHART                        = 0xdf
	DLT_X2E_SERIAL                    = 0xd5
	DLT_X2E_XORAYA                    = 0xd6
	DT_BLK                            = 0x6
	DT_CHR                            = 0x2
	DT_DIR                            = 0x4
	DT_FIFO                           = 0x1
	DT_LNK                            = 0xa
	DT_REG                            = 0x8
	DT_SOCK                           = 0xc
	DT_UNKNOWN                        = 0x0
	DT_WHT                            = 0xe
	ECHO                              = 0x8
	ECHOCTL                           = 0x40
	ECHOE                             = 0x2
	ECHOK                             = 0x4
	ECHOKE                            = 0x1
	ECHONL                            = 0x10
	ECHOPRT                           = 0x20
	EVFILT_AIO                        = -0x3
	EVFILT_FS                         = -0x9
	EVFILT_MACHPORT                   = -0x8
	EVFILT_PROC                       = -0x5
	EVFILT_READ                       = -0x1
	EVFILT_SIGNAL                     = -0x6
	EVFILT_SYSCOUNT                   = 0xe
	EVFILT_THREADMARKER               = 0xe
	EVFILT_TIMER                      = -0x7
	EVFILT_USER                       = -0xa
	EVFILT_VM                         = -0xc
	EVFILT_VNODE                      = -0x4
	EVFILT_WRITE                      = -0x2
	EV_ADD                            = 0x1
	EV_CLEAR                          = 0x20
	EV_DELETE                         = 0x2
	EV_DISABLE                        = 0x8
	EV_DISPATCH                       = 0x80
	EV_ENABLE                         = 0x4
	EV_EOF                            = 0x8000
	EV_ERROR                          = 0x4000
	EV_FLAG0                          = 0x1000
	EV_FLAG1                          = 0x2000
	EV_ONESHOT                        = 0x10
	EV_OOBAND                         = 0x2000
	EV_POLL                           = 0x1000
	EV_RECEIPT                        = 0x40
	EV_SYSFLAGS                       = 0xf000
	EXTA                              = 0x4b00
	EXTB                              = 0x9600
	EXTPROC                           = 0x800
	FD_CLOEXEC                        = 0x1
	FD_SETSIZE                        = 0x400
	FLUSHO                            = 0x800000
	F_ADDFILESIGS                     = 0x3d
	F_ADDSIGS                         = 0x3b
	F_ALLOCATEALL                     = 0x4
	F_ALLOCATECONTIG                  = 0x2
	F_CHKCLEAN                        = 0x29
	F_DUPFD                           = 0x0
	F_DUPFD_CLOEXEC                   = 0x43
	F_FINDSIGS                        = 0x4e
	F_FLUSH_DATA                      = 0x28
	F_FREEZE_FS                       = 0x35
	F_FULLFSYNC                       = 0x33
	F_GETCODEDIR                      = 0x48
	F_GETFD                           = 0x1
	F_GETFL                           = 0x3
	F_GETLK                           = 0x7
	F_GETLKPID                        = 0x42
	F_GETNOSIGPIPE                    = 0x4a
	F_GETOWN                          = 0x5
	F_GETPATH                         = 0x32
	F_GETPATH_MTMINFO                 = 0x47
	F_GETPROTECTIONCLASS              = 0x3f
	F_GETPROTECTIONLEVEL              = 0x4d
	F_GLOBAL