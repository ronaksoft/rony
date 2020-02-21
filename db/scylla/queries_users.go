package scylla

import (
	"github.com/gocql/gocql"
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
	UserMessagesDelete                *queryPool
	UserMessagesInsert                *queryPool
	UserMessagesCountMentions         *queryPool
	UserMessagesCountUnread           *queryPool
	UserMessagesUpdateContentRead     *queryPool
	UserMessagesSelectMessage         *queryPool
	UserMessagesEdit                  *queryPool
	UsersAuthDelete                   *queryPool
	UsersAuthInsert                   *queryPool
	UsersAuthUpdateAccess             *queryPool
	UsersAuthUpdateDevice             *queryPool
	UsersAuthSelect                   *queryPool
	UsersContactsCountAll             *queryPool
	UsersContactsInsert               *queryPool
	UsersContactsGet                  *queryPool
	UsersContactsSelectAll            *queryPool
	UsersDevicesInsert                *queryPool
	UsersDevicesDelete                *queryPool
	UsersDevicesSelect                *queryPool
	UsersDialogsInsert                *queryPool
	UsersDialogsGet                   *queryPool
	UsersDialogsSelect                *queryPool
	UsersDialogsCount                 *queryPool
	UsersDialogsUpdatePin             *queryPool
	UsersDialogsUpdateReadInboxMaxID  *queryPool
	UsersDialogsUpdateReadOutboxMaxID *queryPool
	UsersDraftsDelete                 *queryPool
	UsersDraftsInsert                 *queryPool
	UsersDraftGet                     *queryPool
	UsersSelect                       *queryPool
	UsersUpdateAddBlocked             *queryPool
	UsersUpdateRemoveBlocked          *queryPool
	UsersUpdateAddAuthID              *queryPool
	UsersSignUp                       *queryPool
	UsersUpdateRemoveAuthID           *queryPool
	UsersMarkAsDelete                 *queryPool
	UsersGetByPhone                   *queryPool
	UsersGetByUsername                *queryPool
	UsersSetContactsHash              *queryPool
	UsersSetProfilePhoto              *queryPool
	UsersPhotosAddPhoto               *queryPool
	UsersPhotosDelPhoto               *queryPool
	UsersPhotosSelect                 *queryPool
	UsersUpdatesInsert                *queryPool
	UsersUpdatesGetDiff               *queryPool
	UsersPeersGet                     *queryPool
	UsersPasswordsGet                 *queryPool
	UsersPasswordsSet                 *queryPool
	UsersBlockedInsert                *queryPool
	UsersBlockedDelete                *queryPool
	UsersBlockedSelect                *queryPool
	UsersBlockedCountAll              *queryPool
)

func initUsers() {
	// UsersSelect
	UsersSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt8, names8 := qb.Select(DbTableUsers).
			Where(qb.Eq(CnUserID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt8), names8)
	})

	// UsersSetContactsHash
	UsersSetContactsHash = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableUsers).
			Set(CnContactsHash).
			Where(qb.Eq(CnUserID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersUpdateAddAuthID
	UsersUpdateAddAuthID = newQueryPool(func() *gocqlx.Queryx {
		stmt11, names11 := qb.Update(DbTableUsers).
			Where(qb.Eq(CnUserID)).
			Add(CnAuthIDs).ToCql()
		return gocqlx.Query(Session().Query(stmt11), names11)
	})

	// UsersUpdateRemoveAuthID
	UsersUpdateRemoveAuthID = newQueryPool(func() *gocqlx.Queryx {
		stmt13, names13 := qb.Update(DbTableUsers).
			Where(qb.Eq(CnUserID)).
			Remove(CnAuthIDs).
			ToCql()
		return gocqlx.Query(Session().Query(stmt13), names13)
	})

	// UsersUpdateAddBlocked
	UsersUpdateAddBlocked = newQueryPool(func() *gocqlx.Queryx {
		stmt13, names13 := qb.Update(DbTableUsers).
			Where(qb.Eq(CnUserID)).
			Add(CnBlockedUsers).
			ToCql()
		return gocqlx.Query(Session().Query(stmt13), names13)
	})

	// UsersUpdateRemoveBlocked
	UsersUpdateRemoveBlocked = newQueryPool(func() *gocqlx.Queryx {
		stmt13, names13 := qb.Update(DbTableUsers).
			Where(qb.Eq(CnUserID)).
			Remove(CnBlockedUsers).
			ToCql()
		return gocqlx.Query(Session().Query(stmt13), names13)
	})

	// UsersSetProfilePhoto
	UsersSetProfilePhoto = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableUsers).
			Set(CnPhotoBig, CnPhotoSmall, CnPhotoID).
			Where(qb.Eq(CnUserID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersSignUp
	UsersSignUp = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableUsers).
			Columns(CnUserID, CnFName, CnLName, CnPhone, CnCreatedOn, CnDisabled).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersSignUpBot
	UsersSignUpBot = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableUsers).
			Columns(CnUserID, CnFName, CnUsername, CnAuthIDs, CnIsBot, CnCreatedOn).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersMarkAsDelete
	UsersMarkAsDelete = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableUsers).
			Where(qb.Eq(CnUserID)).
			Set(CnFName, CnLName, CnPhone, CnUsername, CnDeleted).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersGetByUsername
	UsersGetByUsername = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbViewUsersByUsername).
			Where(qb.Eq(CnUsername)).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initUsersAuth() {
	// UsersAuthUpdateDevice
	UsersAuthUpdateDevice = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableUsersAuth).
			Set(CnAppVersion, CnLangCode, CnModel, CnClientID, CnSystemVersion).
			Where(qb.Eq(CnUserID), qb.Eq(CnAuthID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersAuthInsert
	UsersAuthInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableUsersAuth).
			Columns(CnUserID, CnAuthID, CnClientIP, CnCreatedOn).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersAuthDelete
	UsersAuthDelete = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Delete(DbTableUsersAuth).
			Where(qb.Eq(CnUserID), qb.Eq(CnAuthID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersAuthUpdateAccess
	UsersAuthUpdateAccess = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableUsersAuth).
			Set(CnLastAccess, CnClientIP).
			Where(qb.Eq(CnUserID), qb.Eq(CnAuthID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersAuthSelect
	UsersAuthSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersAuth).
			Columns(CnAppVersion, CnSystemVersion, CnLangCode, CnModel, CnAuthID, CnUserID, CnClientIP, CnCreatedOn, CnLastAccess).
			Where(qb.Eq(CnUserID)).Limit(100).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})
}
func initUsersDevices() {
	// UsersDevicesInsert
	UsersDevicesInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableUsersDevices).
			Columns(CnUserID, CnAuthID, CnToken, CnTokenType, CnAppVersion, CnSystemVersion, CnLangCode, CnModel, CnClientID).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)

	})

	// UsersDevicesDelete
	UsersDevicesDelete = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Delete(DbTableUsersDevices).
			Where(qb.Eq(CnUserID), qb.Eq(CnAuthID)).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersDevicesSelect
	UsersDevicesSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersDevices).
			Columns(CnToken, CnTokenType, CnAppVersion, CnSystemVersion, CnLangCode, CnModel, CnAuthID, CnClientID).
			Where(qb.Eq(CnUserID)).Limit(100).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initUsersDialogs() {
	// UsersDialogsInsert
	UsersDialogsInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableUsersDialogs).
			Columns(CnPeerID, CnPeerType, CnUserID, CnTopMessageID, CnTopSenderID).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersDialogsUpdatePin
	UsersDialogsUpdatePin = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableUsersDialogs).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType)).
			Set(CnPinned).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersDialogsUpdateReadOutboxMaxID
	UsersDialogsUpdateReadOutboxMaxID = newQueryPool(func() *gocqlx.Queryx {
		stmt4, names4 := qb.Update(DbTableUsersDialogs).
			Set(CnReadOutboxMaxID).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt4), names4)
	})

	// UsersDialogsUpdateReadInboxMaxID
	UsersDialogsUpdateReadInboxMaxID = newQueryPool(func() *gocqlx.Queryx {
		stmt5, names5 := qb.Update(DbTableUsersDialogs).
			Set(CnReadInboxMaxID).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt5), names5)
	})

	// UsersDialogsGet
	UsersDialogsGet = newQueryPool(func() *gocqlx.Queryx {
		stmt6, names6 := qb.Select(DbTableUsersDialogs).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt6), names6)
	})

	// UsersDialogsSelect
	UsersDialogsSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt6, names6 := qb.Select(DbTableUsersDialogs).
			Where(qb.Eq(CnUserID)).
			Columns(CnPeerType, CnPeerID).
			OrderBy(CnPeerType, gocql.DESC).
			OrderBy(CnPeerID, gocql.DESC).
			ToCql()
		return gocqlx.Query(Session().Query(stmt6), names6)
	})

	// UsersDialogsCount
	UsersDialogsCount = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersDialogs).
			Where(qb.Eq(CnUserID)).
			CountAll().
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initUsersDrafts() {
	// UsersDraftsInsert
	UsersDraftsInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableUsersDrafts).
			Columns(
				CnUserID, CnPeerID, CnPeerType, CnReplyTo,
				CnBody, CnCreatedOn, CnEntities, CnScrambled, CnEditedID,
			).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersDraftsDelete
	UsersDraftsDelete = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Delete(DbTableUsersDrafts).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersDraftsGet
	UsersDraftGet = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersDrafts).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initUsersContacts() {
	// UsersContactsSelectAll
	UsersContactsSelectAll = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersContacts).
			Columns(CnContactID, CnClientID, CnFName, CnLName, CnPhone, CnAccessHash).
			Where(qb.Eq(CnUserID)).
			Limit(2500).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersContactsGet
	UsersContactsGet = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersContacts).
			Columns(CnContactID, CnClientID, CnFName, CnLName, CnPhone, CnAccessHash).
			Where(qb.Eq(CnUserID), qb.Eq(CnContactID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersContactsInsert
	UsersContactsInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableUsersContacts).
			Columns(CnUserID, CnDeviceID, CnContactID, CnClientID, CnFName, CnLName, CnPhone, CnAccessHash).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersContactsCountAll
	UsersContactsCountAll = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersContacts).
			CountAll().
			Where(qb.Eq(CnUserID), qb.Eq(CnContactID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initUsersMessages() {
	// UserMessagesCountUnread
	UserMessagesCountUnread = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersMessages).
			CountAll().
			Where(
				qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType),
				qb.Gt(CnLocalID), qb.LtOrEqNamed(CnLocalID, CnMsgID), qb.Eq(CnInbox),
			).AllowFiltering().ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UserMessagesCountMentions
	UserMessagesCountMentions = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersMessages).
			CountAll().
			Where(
				qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType),
				qb.Gt(CnLocalID), qb.LtOrEqNamed(CnLocalID, CnMsgID), qb.Eq(CnMentioned),
			).AllowFiltering().ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UserMessagesInsert
	UserMessagesInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt2, names2 := qb.Insert(DbTableUsersMessages).
			Columns(
				CnUserID, CnSenderID, CnPeerID, CnPeerType,
				CnSenderMsgID, CnLocalID, CnGlobalID, CnReplyTo,
				CnBody, CnInbox, CnMsgAction, CnMsgActionData,
				CnCreatedOn,
				CnFwdChannelID, CnFwdChannelMsgID, CnFwdSenderID,
				CnEntities, CnMentioned, CnScrambled,
				CnMediaType, CnMediaData, CnReplyMarkup, CnReplyMarkupData,
			).ToCql()
		return gocqlx.Query(Session().Query(stmt2), names2)
	})

	// UserMessagesDelete
	UserMessagesDelete = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Delete(DbTableUsersMessages).
			Where(
				qb.Eq(CnUserID),
				qb.Eq(CnPeerType),
				qb.Eq(CnPeerID),
				qb.In(CnLocalID),
			).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UserMessagesUpdateContentRead
	UserMessagesUpdateContentRead = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableUsersMessages).
			Set(CnContentRead).
			Where(
				qb.Eq(CnUserID),
				qb.Eq(CnPeerType),
				qb.Eq(CnPeerID),
				qb.In(CnLocalID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UserMessagesSelectMessage
	UserMessagesSelectMessage = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersMessages).
			Where(
				qb.Eq(CnUserID),
				qb.Eq(CnPeerType),
				qb.Eq(CnPeerID),
				qb.Eq(CnLocalID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// GetHistoryWithMin
	GetHistoryWithMin = newQueryPool(func() *gocqlx.Queryx {
		cmps := []qb.Cmp{
			qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType), qb.GtOrEqNamed(CnLocalID, CnMinID),
		}
		stmt, names := qb.Select(DbTableUsersMessages).
			Where(cmps...).
			OrderBy(CnPeerType, qb.DESC).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// GetHistoryWithMax
	GetHistoryWithMax = newQueryPool(func() *gocqlx.Queryx {
		cmps := []qb.Cmp{
			qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType), qb.LtOrEqNamed(CnLocalID, CnMaxID),
		}
		stmt, names := qb.Select(DbTableUsersMessages).
			Where(cmps...).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// GetHistoryWithMinMax
	GetHistoryWithMinMax = newQueryPool(func() *gocqlx.Queryx {
		cmps := []qb.Cmp{
			qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType), qb.LtOrEqNamed(CnLocalID, CnMaxID), qb.GtOrEqNamed(CnLocalID, CnMinID),
		}
		stmt, names := qb.Select(DbTableUsersMessages).
			Where(cmps...).
			ToCql()

		return gocqlx.Query(Session().Query(stmt), names)
	})

	// GetHistory
	GetHistory = newQueryPool(func() *gocqlx.Queryx {
		cmps := []qb.Cmp{qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType)}
		stmt, names := qb.Select(DbTableUsersMessages).
			Where(cmps...).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersMessagesEdit
	UserMessagesEdit = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableUsersMessages).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerID), qb.Eq(CnPeerType), qb.Eq(CnLocalID)).
			Set(CnBody, CnEditedOn, CnEntities).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initUsersPhotos() {
	UsersPhotosAddPhoto = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableUsersPhotos).
			Columns(CnUserID, CnPhotoID, CnPhotoBig, CnPhotoSmall).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	UsersPhotosDelPhoto = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Delete(DbTableUsersPhotos).
			Where(qb.Eq(CnUserID), qb.Eq(CnPhotoID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	UsersPhotosSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersPhotos).
			Where(qb.Eq(CnUserID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initUsersUpdates() {
	// UsersUpdatesInsert
	UsersUpdatesInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt7, names7 := qb.Insert(DbTableUsersUpdates).
			Columns(CnUserID, CnUpdateID, CnConstructor, CnObject, CnCreatedOn).ToCql()
		return gocqlx.Query(Session().Query(stmt7), names7)
	})

	UsersUpdatesGetDiff = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersUpdates).
			Columns(CnUpdateID, CnConstructor, CnCreatedOn, CnObject).
			Where(qb.Eq(CnUserID), qb.GtOrEq(CnUpdateID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initBotsPeers() {
	BotsPeersCount = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableBotsPeers).
			CountAll().
			Where(qb.Eq(CnBotID), qb.Eq(CnPeerID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	BotsPeersInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableBotsPeers).
			Columns(CnBotID, CnPeerID).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})
}
func initUsersPasswords() {
	// UsersPasswordsGet
	UsersPasswordsGet = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersPasswords).
			Columns(CnAlgorithm, CnAlgorithmData, CnHint, CnSrpHash, CnSecurityQuestions).
			Where(qb.Eq(CnUserID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersPasswordsSet
	UsersPasswordsSet = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableUsersPasswords).
			Columns(CnUserID, CnAlgorithm, CnAlgorithmData, CnHint, CnSrpHash, CnSecurityQuestions).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initUsersBlocked() {
	// UsersBlockedInsert
	UsersBlockedInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableUsersBlocked).
			Columns(CnUserID, CnPeerID, CnCreatedOn).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersBlockedDelete
	UsersBlockedDelete = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Delete(DbTableUsersBlocked).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersBlockedSelect
	UsersBlockedSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersBlocked).
			Columns(CnPeerID, CnCreatedOn).
			Where(qb.Eq(CnUserID)).
			OrderBy(CnPeerID, qb.ASC).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// UsersBlockedCountAll
	UsersBlockedCountAll = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableUsersBlocked).
			Where(qb.Eq(CnUserID)).
			CountAll().
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
