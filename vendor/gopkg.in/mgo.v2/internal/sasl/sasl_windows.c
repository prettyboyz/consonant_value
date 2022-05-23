#include "sasl_windows.h"

static const LPSTR SSPI_PACKAGE_NAME = "kerberos";

SECURITY_STATUS SEC_ENTRY sspi_acquire_credentials_handle(CredHandle *cred_handle, char *username, char *password, char *domain)
{
	SEC_WINNT_AUTH_IDENTITY auth_identity;
	SECURITY_INTEGER ignored;

	auth_identity.Flags = SEC_WINNT_AUTH_IDENTITY_ANSI;
	auth_identity.User = (LPSTR) username;
	auth_identity.UserLength