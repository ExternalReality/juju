// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"github.com/juju/errors"
	"github.com/juju/mgo"
	"github.com/juju/mgo/txn"
	jujutxn "github.com/juju/txn"

	"github.com/juju/juju/state/storage"
)

// Persistence exposes persistence-layer functionality of State.
type Persistence interface {
	// One populates doc with the document corresponding to the given
	// ID. Missing documents result in errors.NotFound.
	One(collName, id string, doc interface{}) error

	// All populates docs with the list of the documents corresponding
	// to the provided query.
	All(collName string, query, docs interface{}) error

	// Run runs the transaction generated by the provided factory
	// function. It may be retried several times.
	Run(transactions jujutxn.TransactionSource) error

	// NewStorage returns a new blob storage for the environment.
	NewStorage() storage.Storage

	// ApplicationExistsOps returns the operations that verify that the
	// identified application exists.
	ApplicationExistsOps(applicationID string) []txn.Op

	// IncCharmModifiedVersionOps returns the operations necessary to increment
	// the CharmModifiedVersion field for the given application.
	IncCharmModifiedVersionOps(applicationID string) []txn.Op
}

type statePersistence struct {
	st *State
}

// newPersistence builds a new StatePersistence that wraps State.
func (st *State) newPersistence() Persistence {
	return &statePersistence{st: st}
}

// One gets the identified document from the collection.
func (sp statePersistence) One(collName, id string, doc interface{}) error {
	coll, closeColl := sp.st.db().GetCollection(collName)
	defer closeColl()

	err := coll.FindId(id).One(doc)
	if err == mgo.ErrNotFound {
		return errors.NotFoundf(id)
	}
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// All gets all documents from the collection matching the query.
func (sp statePersistence) All(collName string, query, docs interface{}) error {
	coll, closeColl := sp.st.db().GetCollection(collName)
	defer closeColl()

	if err := coll.Find(query).All(docs); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// Run runs the transaction produced by the provided factory function.
func (sp statePersistence) Run(transactions jujutxn.TransactionSource) error {
	if err := sp.st.db().Run(transactions); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// NewStorage returns a new blob storage for the environment.
func (sp *statePersistence) NewStorage() storage.Storage {
	modelUUID := sp.st.ModelUUID()
	// TODO(ericsnow) Copy the session?
	session := sp.st.session
	store := storage.NewStorage(modelUUID, session)
	return store
}

// ApplicationExistsOps returns the operations that verify that the
// identified service exists.
func (sp *statePersistence) ApplicationExistsOps(applicationID string) []txn.Op {
	return []txn.Op{{
		C:      applicationsC,
		Id:     applicationID,
		Assert: isAliveDoc,
	}}
}

// IncCharmModifiedVersionOps returns the operations necessary to increment the
// CharmModifiedVersion field for the given service.
func (sp *statePersistence) IncCharmModifiedVersionOps(applicationID string) []txn.Op {
	return incCharmModifiedVersionOps(applicationID)
}
