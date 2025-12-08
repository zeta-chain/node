package runner

import observertypes "github.com/zeta-chain/node/x/observer/types"

func (r *E2ERunner) RemoveObserver() error {
	observerSet, err := r.ObserverClient.ObserverSet(r.Ctx, &observertypes.QueryObserverSet{})
	if err != nil {
		return err
	}
	r.Logger.Print("🏃 Removing observer from the set : %s", observerSet.Observers[len(observerSet.Observers)-1])
	err = r.ZetaTxServer.RemoveObserver(observerSet.Observers[len(observerSet.Observers)-1])
	return err
}
