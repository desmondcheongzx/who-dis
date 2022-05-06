package pkg

type RRType = uint16

const (
	RR_A     = 1
	RR_NS    = 2
	RR_MD    = 3
	RR_MF    = 4
	RR_CNAME = 5
	RR_SOA   = 6
	RR_MB    = 7
	RR_MG    = 8
	RR_MR    = 9
	RR_NULL  = 10
	RR_WKS   = 11
	RR_PTR   = 12
	RR_HINFO = 13
	RR_MINFO = 14
	RR_MX    = 15
	RR_TXT   = 16
)

type RRClass = uint16

const (
	RR_IN = 1
	RR_CS = 2
	RR_CH = 3
	RR_HS = 4
)

type QType = uint16

const (
	Q_A     = 1
	Q_NS    = 2
	Q_MD    = 3
	Q_MF    = 4
	Q_CNAME = 5
	Q_SOA   = 6
	Q_MB    = 7
	Q_MG    = 8
	Q_MR    = 9
	Q_NULL  = 10
	Q_WKS   = 11
	Q_PTR   = 12
	Q_HINFO = 13
	Q_MINFO = 14
	Q_MX    = 15
	Q_TXT   = 16
	Q_AXFR  = 252
	Q_MAILB = 253
	Q_MAILA = 254
	Q_TSTAR = 255
)

type QClass = uint16

const (
	Q_IN    = 1
	Q_CS    = 2
	Q_CH    = 3
	Q_HS    = 4
	Q_CSTAR = 255
)
