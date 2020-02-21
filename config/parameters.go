package config

import "hash/crc64"

/*
   Creation Time: 2019 - Dec - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var AccessHashTable = crc64.MakeTable(0x23740630002374)
