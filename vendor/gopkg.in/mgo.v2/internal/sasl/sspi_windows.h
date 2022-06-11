// Code adapted from the NodeJS kerberos library:
// 
//   https://github.com/christkv/kerberos/tree/master/lib/win32/kerberos_sspi.h
//
// Under the terms of the Apache License, Version 2.0:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
#ifndef SSPI_WINDOWS_H
#define SSPI_WINDOWS_H

#define SECURITY_WIN32 1

#include <windows.h>
#include <sspi.h>

int load_secur32_dll();

SECURITY_STATUS SEC_ENTRY call_sspi_encrypt_message(PCtxtHandle phContext, unsigned long fQOP, PSecBufferDesc pMessage, unsigned long MessageSeqNo);