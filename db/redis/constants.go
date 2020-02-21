package redis

/*
   Creation Time: 2019 - Sep - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// REDIS KEY NAMES
const (
	RkRsaFingerPrint        = "RSA_FP" // RSA Finger Print
	RkPQ                    = "PQ"
	RkDhFingerPrint         = "DHFP" // Diffie-Helman Finger Print
	RkPhoneVerification     = "PHONE_VERIFY"
	RkPhoneVerificationCode = "CODE"
	RkPhoneVerificationTry  = "TRY"
	RkDialogs               = "DLGS"
	RkCounter               = "CNT"
	RkUnread                = "UNRD"
	RkMentioned             = "MNT"
	RkPinnedDialogs         = "PDLGS"
	RkSystemInfo            = "SYS_INFO"
	RkFilePartSize          = "FILE_PS"
	RkSrpB                  = "SRP_B"
	RkTempBind              = "TMP_B"
)
