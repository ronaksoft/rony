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
	LabelSet                    *queryPool
	LabelGet                    *queryPool
	LabelUpdateCount            *queryPool
	LabelDelete                 *queryPool
	LabelGetByUser              *queryPool
	LabelGetItems               *queryPool
	LabelCountItems             *queryPool
	LabelGetItemsByMsgIDWithMin *queryPool
	LabelGetItemsByMsgIDWithMax *queryPool
	LabelGetItemsByMsgID        *queryPool
	LabelInsertItem             *queryPool
	LabelDeleteItem             *queryPool
	LabelRemoveFromUserMessage  *queryPool
	LabelAddToUserMessage       *queryPool
)

func initLabels() {
	LabelSet = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableLabels).
			Columns(CnUserID, CnLabelID, CnName, CnColour).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelGet = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableLabels).
			Where(qb.Eq(CnUserID), qb.In(CnLabelID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelUpdateCount = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableLabels).
			Where(qb.Eq(CnUserID), qb.Eq(CnLabelID)).
			Set(CnCount).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelDelete = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Delete(DbTableLabels).
			Where(qb.Eq(CnUserID), qb.Eq(CnLabelID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelGetByUser = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableLabels).
			Where(qb.Eq(CnUserID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelGetItems = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableLabelsItems).
			Columns(CnPeerType, CnPeerID, CnMsgID).
			Where(qb.Eq(CnUserID), qb.Eq(CnLabelID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelGetItemsByMsgID = newQueryPool(func() *gocqlx.Queryx {
		cmps := []qb.Cmp{
			qb.Eq(CnUserID), qb.Eq(CnLabelID),
		}
		stmt, names := qb.Select(DbViewLabelsItemsByMsgID).
			Columns(CnPeerType, CnPeerID, CnMsgID).
			OrderBy(CnLabelID, qb.DESC).
			OrderBy(CnMsgID, qb.DESC).
			Where(cmps...).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelGetItemsByMsgIDWithMin = newQueryPool(func() *gocqlx.Queryx {
		cmps := []qb.Cmp{
			qb.Eq(CnUserID), qb.Eq(CnLabelID), qb.GtOrEqNamed(CnMsgID, CnMinID),
		}
		stmt, names := qb.Select(DbViewLabelsItemsByMsgID).
			Columns(CnPeerType, CnPeerID, CnMsgID).
			Where(cmps...).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})
	LabelGetItemsByMsgIDWithMax = newQueryPool(func() *gocqlx.Queryx {
		cmps := []qb.Cmp{
			qb.Eq(CnUserID), qb.Eq(CnLabelID), qb.LtOrEqNamed(CnMsgID, CnMaxID),
		}
		stmt, names := qb.Select(DbViewLabelsItemsByMsgID).
			Columns(CnPeerType, CnPeerID, CnMsgID).
			OrderBy(CnLabelID, qb.DESC).
			OrderBy(CnMsgID, qb.DESC).
			Where(cmps...).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelInsertItem = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableLabelsItems).
			Columns(CnUserID, CnLabelID, CnPeerType, CnPeerID, CnMsgID, CnCreatedOn).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelDeleteItem = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Delete(DbTableLabelsItems).
			Where(qb.Eq(CnUserID), qb.Eq(CnLabelID), qb.Eq(CnPeerType), qb.Eq(CnPeerID), qb.Eq(CnMsgID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelAddToUserMessage = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableUsersMessages).
			Add(CnLabelIDs).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerType), qb.Eq(CnPeerID), qb.Eq(CnLocalID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelRemoveFromUserMessage = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Update(DbTableUsersMessages).
			Remove(CnLabelIDs).
			Where(qb.Eq(CnUserID), qb.Eq(CnPeerType), qb.Eq(CnPeerID), qb.Eq(CnLocalID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	LabelCountItems = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableLabelsItems).
			CountAll().
			Where(qb.Eq(CnUserID), qb.Eq(CnLabelID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})
}
