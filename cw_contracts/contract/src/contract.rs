use cosmwasm_std::{Uint256, Uint64};
use cosmwasm_std::{entry_point, to_binary};
use cosmwasm_std::{Deps, DepsMut, Env, MessageInfo};
use cosmwasm_std::{QueryResponse, Response, StdError, StdResult};

use schemars::JsonSchema;
use thiserror::Error;

use serde::{Deserialize, Serialize};

use zeta_std::{WatchlistQueryResponse,ZetaCoreMsg,ZetaCoreQuery};

#[derive(Error, Debug)]
pub enum WatcherError {
    #[error("{0}")]
    Std(#[from] StdError),
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Instantiate
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)] //JsonSchema removed
pub struct InstantiateMsg {}

#[entry_point]
pub fn instantiate(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, WatcherError> {
    Ok(Response::default())
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Execute
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ExecuteMsg {
    AddToWatchList {
        chain: String,
        nonce: Uint64,
        tx_hash: String,
    }
}

#[entry_point]
pub fn execute(
    deps: DepsMut<ZetaCoreQuery>,
    _env: Env,
    _info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response<ZetaCoreMsg>, WatcherError> {
    match msg {
        ExecuteMsg::AddToWatchList { chain,nonce,tx_hash } => {
            let add_watchlist_msg = SifchainMsg::Swap {
                sent_asset: "rowan".to_string(),
                received_asset: "ceth".to_string(),
                sent_amount: amount.to_string(),
                min_received_amount: "0".to_string(),
            };

            Ok(Response::new()
                .add_attribute("action", "swap")
                .add_message(swap_msg))
        }
    }
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Query
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum QueryMsg {
    Pool { external_asset: String },
}

#[entry_point]
pub fn query(deps: Deps<SifchainQuery>, _env: Env, msg: QueryMsg) -> StdResult<QueryResponse> {
    match msg {
        QueryMsg::Pool { external_asset } => to_binary(&query_pool(deps, external_asset)?),
    }
}

fn query_pool(deps: Deps<SifchainQuery>, external_asset: String) -> StdResult<PoolResponse> {
    let req = SifchainQuery::Pool { external_asset }.into();
    deps.querier.query(&req)
}
