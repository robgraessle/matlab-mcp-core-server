// Copyright 2025 The MathWorks, Inc.

package globalmatlab

import (
	"context"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type MATLABManager interface {
	StartMATLABSession(ctx context.Context, sessionLogger entities.Logger, startRequest entities.SessionDetails) (entities.SessionID, error)
	GetMATLABSessionClient(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) (entities.MATLABSessionClient, error)
}

type MATLABRootSelector interface {
	SelectFirstMATLABVersionOnPath(ctx context.Context, logger entities.Logger) (string, error)
}

type MATLABStartingDirSelector interface {
	SelectMatlabStartingDir() (string, error)
}

type GlobalMATLAB struct {
	matlabManager             MATLABManager
	matlabRootSelector        MATLABRootSelector
	matlabStartingDirSelector MATLABStartingDirSelector

	lock              *sync.Mutex
	initializeOnce    *sync.Once
	matlabRoot        string
	matlabStartingDir string
	sessionID         entities.SessionID
	cachedStartupErr  error
}

func New(
	matlabManager MATLABManager,
	matlabRootSelector MATLABRootSelector,
	matlabStartingDirSelector MATLABStartingDirSelector,
) *GlobalMATLAB {
	return &GlobalMATLAB{
		matlabManager:             matlabManager,
		matlabRootSelector:        matlabRootSelector,
		matlabStartingDirSelector: matlabStartingDirSelector,

		lock:           &sync.Mutex{},
		initializeOnce: &sync.Once{},
	}
}

func (g *GlobalMATLAB) Client(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.initializeOnce.Do(func() {
		err := g.initializeMATLABStartupVariables(ctx, logger)
		if err != nil {
			g.cachedStartupErr = err
		}
	})

	if g.cachedStartupErr != nil {
		return nil, g.cachedStartupErr
	}

	if err := g.ensureMATLABSessionIsReady(ctx, logger); err != nil {
		g.cachedStartupErr = err
		return nil, err
	}

	client, err := g.matlabManager.GetMATLABSessionClient(ctx, logger, g.sessionID)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (g *GlobalMATLAB) ensureMATLABSessionIsReady(ctx context.Context, logger entities.Logger) error {

	var sessionIDZeroValue entities.SessionID
	if g.sessionID == sessionIDZeroValue {
		sessionID, err := g.matlabManager.StartMATLABSession(ctx, logger, entities.LocalSessionDetails{
			MATLABRoot:             g.matlabRoot,
			IsStartingDirectorySet: g.matlabStartingDir != "",
			StartingDirectory:      g.matlabStartingDir,
			ShowMATLABDesktop:      true,
		})
		if err != nil {
			return err
		}

		g.sessionID = sessionID
	}

	return nil
}

func (g *GlobalMATLAB) initializeMATLABStartupVariables(ctx context.Context, logger entities.Logger) error {
	matlabRoot, err := g.matlabRootSelector.SelectFirstMATLABVersionOnPath(ctx, logger)
	if err != nil {
		return err
	}

	g.matlabRoot = matlabRoot

	matlabStartingDirectory, err := g.matlabStartingDirSelector.SelectMatlabStartingDir()
	if err != nil {
		logger.WithError(err).Warn("failed to determine MATLAB starting directory, proceeding without one")
		return nil
	}

	g.matlabStartingDir = matlabStartingDirectory
	return nil
}
