package scylla

import (
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

/*
   Creation Time: 2019 - Sep - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	DocumentsInsert       *queryPool
	DocumentsSelect       *queryPool
	FilesTempDelete       *queryPool
	FilesTempInsert       *queryPool
	FilesTempSelect       *queryPool
	GroupsSelect          *queryPool
	GroupsSetProfilePhoto *queryPool
	GroupsPhotosAddPhoto  *queryPool
	GroupsPhotosDelPhoto  *queryPool
	MessagesGetByGlobalID *queryPool
	PhonesUsersInsert     *queryPool
	PhonesUsersSelect     *queryPool
	SystemInstancesInsert *queryPool
	UsersSignUpBot        *queryPool
	BotsPeersCount        *queryPool
	BotsPeersInsert       *queryPool
	GetHistory            *queryPool
	GetHistoryWithMax     *queryPool
	GetHistoryWithMin     *queryPool
	GetHistoryWithMinMax  *queryPool
	DocumentCounterInc    *queryPool
	DocumentCounterGet    *queryPool
)

func initDbQueries() {
	initAuth()
	initAuthKey()
	initDocuments()
	initFilesTemp()
	initGroups()
	initMessages()
	initPhones()
	initPhonesUsers()
	initUsers()
	initUsersAuth()
	initUsersBlocked()
	initUsersDevices()
	initUsersDialogs()
	initUsersDrafts()
	initUsersPhotos()
	initUsersContacts()
	initUsersMessages()
	initUsersUpdates()
	initUsersPasswords()
	initBotsPeers()
	initLabels()

	SystemInstancesInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableSystemInstances).
			Where(qb.Eq(CnInstanceID)).
			Set(CnGlobalMsgCounter).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	UsersPeersGet = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersPeers).
			Columns(CnNotifyFlags, CnNotifyMuteUntil, CnNotifySound).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerType), qb.Eq(CnPeerID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	DocumentCounterInc = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableDocumentsCounters).
			Add(CnRefCount).
			Where(qb.Eq(CnClusterID), qb.Eq(CnFileID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names).Idempotent(false)
	})

	DocumentCounterGet = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableDocumentsCounters).Columns(CnRefCount).
			Where(qb.Eq(CnClusterID), qb.Eq(CnFileID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}

func initDocuments() {
	// DocumentsInsert
	DocumentsInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableDocuments).
			Columns(CnClusterID, CnFileID, CnMimeType, CnFileSize, CnAttributes, CnAccessHash, CnUserID, CnThumb, CnMD5Checksum).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// DocumentsSelect
	DocumentsSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableDocuments).
			Where(qb.Eq(CnClusterID), qb.Eq(CnFileID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initFilesTemp() {
	// FilesTempInsert
	FilesTempInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableFilesTemp).
			Columns(CnClusterID, CnUserID, CnFileID, CnFileSize, CnDocumentID, CnAccessHash, CnMD5Checksum).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// FilesTempSelect
	FilesTempSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableFilesTemp).
			Where(qb.Eq(CnUserID), qb.Eq(CnFileID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// FilesTempDelete
	FilesTempDelete = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Delete(DbTableFilesTemp).
			Where(qb.Eq(CnUserID), qb.Eq(CnFileID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})
}
func initGroups() {
	GroupsSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableGroups).
			Where(qb.Eq(CnGroupID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	GroupsSetProfilePhoto = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableGroups).
			Set(CnPhotoBig, CnPhotoSmall, CnPhotoID).
			Where(qb.Eq(CnGroupID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	GroupsPhotosAddPhoto = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableGroupsPhotos).
			Columns(CnGroupID, CnPhotoID, CnPhotoBig, CnPhotoSmall).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	GroupsPhotosDelPhoto = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Delete(DbTableGroupsPhotos).
			Where(qb.Eq(CnGroupID), qb.Eq(CnPhotoID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initMessages() {
	// MessagesGetByGlobalID
	MessagesGetByGlobalID = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbViewUserMessagesByGlobalID).
			Where(qb.Eq(CnGlobalID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})
}
func initPhones() {
	UsersGetByPhone = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbViewUsersByPhone).
			Columns(CnUserID).
			Where(qb.Eq(CnPhone)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})
}
func initPhonesUsers() {
	// PhonesUsersInsert
	PhonesUsersInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTablePhonesUsers).
			Columns(CnPhone, CnUserID, CnClientID, CnFName, CnLName, CnScrambled).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// PhonesUsersInsert
	PhonesUsersSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTablePhonesUsers).
			Where(qb.Eq(CnPhone)).Limit(250).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
